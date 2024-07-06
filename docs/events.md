# Events

* Win - event of winning the auction and sends the ad to the user.
* Impression - sends the ad to the user browser, but the user does not see it yet.
* View - the user sees the ad.
* Click - the user clicks on the ad.
* Lead - the user performs the target action.
* Direct - display popup or redirect to the advertiser's site. (it's the same as a click but special for Popup ads)

## Event flow

1. Run collecting of advertisements for the request by user.
2. Run the auction.
3. Sends Win and Impression events for all selected advertisements.
4. Response to the user with the selected advertisement and lists of postback events to count the user actions (view, click).
5. Count the user actions (view, click, direct) by postback events for user actions.
6. Count the user actions (lead) by postback events for user actions.
