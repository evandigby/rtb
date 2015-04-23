package redis

import (
	"github.com/evandigby/rtb"
	"github.com/fzzy/radix/extra/pool"
	"github.com/fzzy/radix/redis"
	"strings"
	"time"
)

type RedisDataAccess struct {
	appDomain string
	network   string
	addr      string

	debitIfNotZeroSha string
	dailyBudgetSha    string

	zunionOutputKey string

	pool *pool.Pool
}

func (da *RedisDataAccess) withDomain(accountKey string) string {
	return da.appDomain + ":" + accountKey
}

func (da *RedisDataAccess) DeleteKeys(keys []string) {
	client, err := da.pool.Get()

	if err != nil {
		panic(err)
	}

	defer da.pool.CarefullyPut(client, &err)

	keysToDelete := make([]interface{}, len(keys))

	for i, val := range keys {
		keysToDelete[i] = da.withDomain(val)
	}

	client.Cmd("DEL", keysToDelete)

	err = client.GetReply().Err
}

func (da *RedisDataAccess) HGetInt64(accountKey string, key string) int64 {
	client, err := da.pool.Get()

	if err != nil {
		panic(err)
	}

	defer da.pool.CarefullyPut(client, &err)

	reply, err := client.Cmd("HGET", da.withDomain(accountKey), key).Int64()

	if err != nil {
		panic(err)
	}

	return reply
}

func (da *RedisDataAccess) HSetInt64(accountKey string, key string, val int64) {
	client, err := da.pool.Get()

	if err != nil {
		panic(err)
	}

	defer da.pool.CarefullyPut(client, &err)

	reply := client.Cmd("HSET", da.withDomain(accountKey), key, val)

	if reply.Err != nil {
		panic(reply.Err)
	}
}

func (da *RedisDataAccess) GetInt64(accountKey string) (success bool, result int64) {
	client, err := da.pool.Get()

	if err != nil {
		panic(err)
	}

	defer da.pool.CarefullyPut(client, &err)

	reply := client.Cmd("GET", da.withDomain(accountKey))

	if reply.Err != nil {
		panic(reply.Err)
	}

	if reply.Type == redis.NilReply {
		return false, 0
	} else {
		val, err := reply.Int64()

		if err != nil {
			panic(err)
		}

		return true, val
	}
}

func (da *RedisDataAccess) SetInt64(accountKey string, val int64) {
	client, err := da.pool.Get()

	if err != nil {
		panic(err)
	}

	defer da.pool.CarefullyPut(client, &err)
	err = client.Cmd("SET", da.withDomain(accountKey), val).Err

	if err != nil {
		panic(err)
	}
}

func (da *RedisDataAccess) GetSetMembers(setKey string) []string {
	client, err := da.pool.Get()

	if err != nil {
		panic(err)
	}

	defer da.pool.CarefullyPut(client, &err)
	reply, err := client.Cmd("SMEMBERS", da.withDomain(setKey)).List()

	if err != nil {
		panic(err)
	}

	return reply
}

func (da *RedisDataAccess) AddMembersToSet(setKey string, members []interface{}) {
	client, err := da.pool.Get()

	if err != nil {
		panic(err)
	}

	defer da.pool.CarefullyPut(client, &err)
	values := make([]interface{}, 0, len(members)+1)

	values = append(values, da.withDomain(setKey))

	values = append(values, members)
	reply := client.Cmd("SADD", values)

	if reply.Err != nil {
		panic(reply.Err)
	}

	//PipelineQueueEmptyError
	for reply.Err == nil {
		reply = client.GetReply()
	}
}

func (da *RedisDataAccess) AddMembersToSortedSets(keyValues map[string]SortedSetMember) {
	client, err := da.pool.Get()

	if err != nil {
		panic(err)
	}

	defer da.pool.CarefullyPut(client, &err)
	for key, value := range keyValues {
		client.Append("ZADD", da.withDomain(key), value.Score, value.Member)
	}

	reply := client.GetReply()

	if reply.Err != nil {
		panic(reply.Err)
	}
}

func (da *RedisDataAccess) SortedSetUnion(keys []string) []string {
	client, err := da.pool.Get()

	if err != nil {
		panic(err)
	}

	defer da.pool.CarefullyPut(client, &err)
	zunionArguments := make([]interface{}, 0, len(keys)+4)

	zunionArguments = append(zunionArguments, da.zunionOutputKey)
	zunionArguments = append(zunionArguments, len(keys))

	for _, key := range keys {
		zunionArguments = append(zunionArguments, da.withDomain(key))
	}

	zunionArguments = append(zunionArguments, "AGGREGATE")
	zunionArguments = append(zunionArguments, "MAX")

	client.Append("MULTI")
	client.Append("ZUNIONSTORE", zunionArguments)
	client.Append("ZREVRANGE", da.zunionOutputKey, "0", "-1")
	client.Append("EXEC")

	reply := client.GetReply()
	reply = client.GetReply()
	reply = client.GetReply()
	reply = client.GetReply()

	if len(reply.Elems) < 2 {
		panic("Did not get a list back")
	}
	result, err := reply.Elems[1].List()

	if err != nil {
		panic(err)
	}

	return result
}

// KEYS[1] is the account key
// ARGV[1] is the daily budget
// ARGV[2] is the expiration time
func DailyBudgetScript() string {
	return `
	if redis.call('EXISTS', KEYS[1]) ~= 1 then
		redis.call('SET', KEYS[1], ARGV[1])
		return redis.call('EXPIREAT', KEYS[1], ARGV[2])
	end`
}

// KEYS[1] is the account key
// ARGV[1] is the amount to debit
func DebitIfNotZeroScript() string {
	return `
	if tonumber(redis.call('GET', KEYS[1])) >= tonumber(ARGV[1]) then
		return redis.call('DECRBY', KEYS[1], ARGV[1])
	else
		return nil
	end`
}

func (da *RedisDataAccess) getDebitIfNotZeroSha() string {
	client, err := da.pool.Get()

	if err != nil {
		panic(err)
	}

	defer da.pool.CarefullyPut(client, &err)
	if da.dailyBudgetSha == "" {
		sha, err := client.Cmd("SCRIPT", "LOAD", DebitIfNotZeroScript()).Str()

		if err != nil {
			panic(err)
		}

		da.dailyBudgetSha = sha
	}

	return da.dailyBudgetSha
}

func (da *RedisDataAccess) getDailyBudgetSha() string {
	client, err := da.pool.Get()

	if err != nil {
		panic(err)
	}

	defer da.pool.CarefullyPut(client, &err)
	if da.debitIfNotZeroSha == "" {
		sha, err := client.Cmd("SCRIPT", "LOAD", DailyBudgetScript()).Str()

		if err != nil {
			panic(err)
		}

		da.debitIfNotZeroSha = sha
	}

	return da.debitIfNotZeroSha
}

func (da *RedisDataAccess) ExpireKey(key string, expirationTime time.Time) {
	client, err := da.pool.Get()

	if err != nil {
		panic(err)
	}

	defer da.pool.CarefullyPut(client, &err)
	reply := client.Cmd("EXPIREAT", da.withDomain(key), expirationTime.Unix())

	if reply.Err != nil {
		panic(reply.Err)
	}
}

func (da *RedisDataAccess) DebitIfNotZero(accountKey string, amount int64, dailyBudget int64, dailyBudgetExpiration time.Time) (remainingDailyBudgetInMicroCents int64, err error) {
	client, err := da.pool.Get()

	if err != nil {
		panic(err)
	}

	defer da.pool.CarefullyPut(client, &err)
	accountKeyWithDomain := da.withDomain(accountKey)
	debitIfNotZeroSha := da.getDebitIfNotZeroSha()
	dailyBudgetSha := da.getDailyBudgetSha()

	client.Append("MULTI")
	client.Append("EVALSHA", dailyBudgetSha, 1, accountKeyWithDomain, dailyBudget, dailyBudgetExpiration.Unix())
	client.Append("EVALSHA", debitIfNotZeroSha, 1, accountKeyWithDomain, amount)
	client.Append("GET", accountKeyWithDomain)
	client.Append("EXEC")

	client.GetReply()
	client.GetReply()
	client.GetReply()
	client.GetReply()
	transaction := client.GetReply()
	if transaction.Err != nil {
		return 0, rtb.NewTransactionError(transaction.Err.Error(), false)
	}

	if err != nil {
		return 0, rtb.NewTransactionError(err.Error(), false)
	}

	amountValid := transaction.Elems[2].Type != redis.NilReply && transaction.Elems[2].Err == nil

	if amountValid {
		remainingDailyBudgetInMicroCents, _ = transaction.Elems[2].Int64()
	}
	if transaction.Elems[1].Type == redis.NilReply {
		return remainingDailyBudgetInMicroCents, rtb.NewTransactionError("Insufficient daily funds.", amountValid)
	} else {
		return remainingDailyBudgetInMicroCents, nil
	}
}

// Should only be used for testing
func (da *RedisDataAccess) GetKeys() []string {
	client, err := da.pool.Get()

	if err != nil {
		panic(err)
	}

	defer da.pool.CarefullyPut(client, &err)
	reply, err := client.Cmd("KEYS", da.withDomain("*")).List()

	if err != nil {
		panic(err)
	}

	withoutDomain := make([]string, len(reply))
	prefix := da.appDomain + ":"
	for i, val := range reply {
		withoutDomain[i] = strings.TrimPrefix(val, prefix)
	}

	return withoutDomain
}

func NewRedisDataAccess(network string, addr string, appDomain string, poolSize int) NoDbDataAccess {
	da := new(RedisDataAccess)
	da.network = network
	da.addr = addr
	da.appDomain = appDomain

	da.zunionOutputKey = da.withDomain("unionoutput")
	pool, err := pool.NewPool(network, addr, poolSize)

	if err != nil {
		panic(err)
	}
	da.pool = pool
	return da
}
