/*

# HTTP headers

@link https://en.wikipedia.org/wiki/List_of_HTTP_header_fields

In HTTP protocol exists special time header marks with the time of the message was originated and lifetime.

Date - The date and time that the message was originated (in "HTTP-date" format as defined by [https://tools.ietf.org/html/rfc7231#section-7.1.1.1](RFC 7231 Date/Time Formats)). ```Date: Tue, 15 Nov 1994 08:12:31 GMT```.
Expires - Gives the date/time after which the response is considered stale (in "HTTP-date" format as defined by [https://tools.ietf.org/html/rfc7231](RFC 7231)). ```Expires: Thu, 01 Dec 1994 16:00:00 GMT```

However, all this time marks restricted in seconds what is not acceptable in RTB integration because of time counting in Milliseconds, and maximal execute time usually in Milliseconds as well.
That's why we have to create new custom HTTP header with the message oridinate time in Milliseconds.

X-Request-Ts: 1549107735603
X-Response-Accepted-Ts: 1549108715201
X-Response-Ts: 1549108715201

*/
package openlatency
