# rtb Real Time Bidding Interface
Proof of concept real time bidding library. This library can be tested using the p.o.c server implementation at https://github.com/evandigby/rtbhost. 

## Implementation Notes

#### Target Matching
- The redis server is used to match campaigns to targets by storing sets that point in "both directions": Sets of all targets a campaign requires, and sets of all campaigns for a specific target.
- When a request comes in, it compiles a list of targets for that request, and then does a union on the sets associated with each target. This produces a list of campaign Ids that contain *ANY* of the targets (an OR relationship).
- There is still a to-do item to implement the *AND* relationship for targeting, allowing campaigns to require all of their targets, or combinations of their targets, not just any of the targets.

#### Bid Decisions
- The target matching sets in redis are stored as sorted sets whos score is their bid CPM. 
- When the set union is returned from target matching, it's already sorted from highest bid CPM to lowest.
- In most cases, the bidder will pick the first one. If the first one does not have available budget, the bidder will move down the list.
- Pacing can be implemented using the "BidPacer" interface. Every time a bidder matches a campaign, it will first ask that campaign if it "can bid" through the pacer.
- There is a sample "time segmented" pacer that will divide the remaining daily budget over the remaining time in the day, and break that into chunks of a specified time segment length. It will not allow any campaign to bid that has exceeded its number of bids for that time segment. This essentially granulates remaining daily budget into smaller chunks. 

#### Bid Responses
- The system will respond with a 200 HTTP Status Code for a bid, and a 204 HTTP Status Code for a no-bid. 
- Right now bid responses are largely stubbed out. These of course need to be completed for any production implementation.
- It's also worth implementing no-bid reasons on 204.

#### Remaining Daily Spending Budget
- The redis server is used as a quick way to cache remaining daily budgets to allow multiple instances of a host using this library to coordinate over the network. 

#### Transaction logging
- The redis server is *NOT* designed to act as a reliable transaction log. 
- Any production implementation of this real time bidder should implement a bomb proof transaction log to maintain accurate accounting records.
- It would also be prudent for another system to consume the transaction log, either directly from the bidder, or through the transaction logger, and periodically audit and update the values stored in the redis instance

#### Error handling
- The current implementation should be stable within the context of what is implemented, but will panic at any errors that arise from unforseen circumstances. This should be fleshed out to provide more detailed error logging and reporting for better debugging in a production environment. 

#### Testing
- You will find tests in the root, inmemory, and redis folders. 
- The root tests are simple unit tests for core functions, such as money conversions.
- inmemory tests are pure unit tests of the bidder using mocks (see rtb/mocks) to stub out any data access
- redis tests are integration tests of the redis components. 

## System Requirements
The rtbhost requires the following:
- Go
- Access to a redis instance.
- Access to an amqp server (if you want to utilize logging).

## Usage
Regarding usage:
- Build rtbhost in the developer home directory:
	go build github.com/evandigby/rtbhost 
- Reset the data store in redis by using the shell script, found in the developer home directory:
	./clear
- Run the host on the default port by executing 
	./rtbhost
- See the many usage options by executing
	./rtbhost --help
- Launch two processes running at ports 8000 and 8001 using the shell script “launch”:
	./launch
- Any parameters passed into launch will be passed onto both rtbhost instances.
- By default, the rtbhost will output a status update on both campaigns every 5 seconds. The shell script will wait 2 seconds before starting the second process to ensure the logging doesn’t overlap.
- You can use the “logverbose” option to force rtbhost to log every request (to stdout by default):
	./launch --logverbose
- There are many other command line options. Please feel free to explore them!
- You can run tests test using:
	go test github.com/evandigby/rtb
	go test github.com/evandigby/rtb/inmemory
	go test github.com/evandigby/rtb/redis 
