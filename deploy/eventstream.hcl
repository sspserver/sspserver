variable "fmt" {
  datetime = "2016-01-02 15:04:05"
}

stores {
  statistic {
    connect = "@env:ES_STORE_STATISTIC_DB_CONNECT"
    // buffer  = 1000
  }
}

// Source could be any supported stream service like kafka, nats, etc...
sources {
  events {
    connect = "@env:ES_SOURCE_EVENTS_CONNECT"
    format  = "json"
  }

  rtb_wins {
    connect = "@env:ES_SOURCE_WINS_CONNECT"
    format  = "json"
  }

  user_info {
    connect = "@env:ES_SOURCE_USERINFO_CONNECT"
    format  = "json"
  }
}

// Streams it's pipelines which have source and destination store
streams {
  events {
    store   = "statistic"
    source  = "events"
    target  = "stats.events_local"
    fields  = [
      "timemark=tm:unixnano",               // DateTime
      "delay=dl:uint",                      // UInt64
      "duration=d:uint",                    // UInt64
      "service=srv:fix*16",                 // FixedString(16)
      "cluster=cl:fix*2",                   // FixedString(2)
      "event=e",                            // Event Type
      "status=st:uint8",                    // Status: 0 - undefined, 1 - success, 2 - failed, 3 - compromised
      // Accounts link information
      "project=pr:uint",                    // UInt64
      "pub_company=pcb:uint",               // UInt64
      "adv_company=acv:uint",               // UInt64
      // Source
      "aucid=auc:uuid",                     // FixedString(16)  -- Internal Auction ID
      "auctype=auctype:uint8",              // Aution type 1 - First price, 2 - Second price
      "impid=imp:uuid",                     // FixedString(16)
      "impadid=impad:uuid",                 // FixedString(16)
      "extaucid=eauc",                      // RTB Request/Response ID
      "extimpid=eimp",                      // RTB Imp ID
      "extzoneid=extz",                     // RTB Zone ID (tagid)
      "source=sid:uint",                    // UInt64
      "network=net",                        // String
      "access_point=acp:uint",              // UInt64
      // State Location
      "platform=pl:uint8",                  // UInt8
      "domain=dm",                          // String
      "app:int",                            // UInt64
      "zone=z:int",                         // UInt64
      "pixel=pxl:uint",                     // UInt64
      "campaign=cmp:int",                   // UInt64
      "format=fmt:uint32",                  // UInt32
      "ad=ad:uint",                         // UInt64
      "ad_w=aw:uint32",                     // UInt32
      "ad_h=ah:uint32",                     // UInt32
      "src_url=su",                         // String
      "win_url=wu",                         // String
      "url=u",                              // String
      "jumper=j:int",                       // UInt64
      // Money section
      "pricing_model=pm:uint8",             // UInt8
      "purchase_view_price=pvpr:uint",      // UInt64
      "purchase_click_price=pcpr:uint",     // UInt64
      "purchase_lead_price=plpr:uint",      // UInt64
      "potential_view_price=ptvpr:uint",    // UInt64
      "potential_click_price=ptcpr:uint",   // UInt64
      "potential_lead_price=ptlpr:uint",    // UInt64
      "view_price=vpr:uint",                // UInt64
      "click_price=cpr:uint",               // UInt64
      "lead_price=lpr:uint",                // UInt64
      "competitor=cmid:uint",               // UInt64
      "competitor_source=cmsrc:uint",       // UInt64
      "competitor_ecpm=cmecpm:int",         // UInt64
      // User IDENTITY
      "udid=udi",                           // FixedString(16)
      "uuid=uui:uuid",                      // FixedString(16)
      "sessid=ses:uuid",                    // FixedString(16)
      "fingerprint=fpr:uuid",               // String
      "etag=etg",                           // String
      // Targeting
      "carrier=car:uint",                   // UInt64
      "country=cc:fix*2",                   // FixedString(2)
      "city=ct:fix*5",                      // FixedString(5)
      "latitude=lt:float",                  // Float64
      "longitude=lg:float",                 // Float64
      "language=lng:fix*5",                 // FixedString(5)
      "ip:ip",                              // IPv6
      "ref",                                // String
      "page_url=page",                      // String
      "ua",                                 // String
      "device_type=dvt:uint32",             // UInt32
      "device=dv:uint32",                   // UInt32
      "os:uint32",                          // UInt32
      "browser=br:uint32",                  // UInt32
      "categories=c:[]int32",               // Array(Int32)
      "adblock=ab:uint8",                   // UInt8
      "private=prv:uint8",                  // UInt8
      "robot=rt:uint8",                     // UInt8
      "proxy=pt:uint8",                     // UInt8
      "backup=bt:uint8",                    // UInt8
      "x:int32",                            // Int32
      "y:int32",                            // Int32
      "w:int32",                            // Int32
      "h:int32",                            // Int32

      "subid1=sd1",
      "subid2=sd2",
      "subid3=sd3",
      "subid4=sd4",
      "subid5=sd5",
    ]
    metrics = [
      {
        name = "event.counter"
        type = "counter"
        tags {
          action   = "{{e}}"
          language = "{{lng}}"
        }
      }
    ]
  }

  rtb_wins {
    store   = "statistic"
    source  = "rtb_wins"
    target  = "stats.rtb_wins"
    fields  = [
      "timemark=tm:unixnano",               // DateTime
      "delay=dl:uint",                      // UInt64
      "duration=d:uint",                    // UInt64
      "service=srv:fix*16",                 // FixedString(16)
      "cluster=cl:fix*2",                   // FixedString(2)
      "aucid=auc:uuid",                     // FixedString(16)  -- Internal Auction ID
      "source=sid:uint",                    // UInt64
      "network=net",                        // String
      "access_point=acp:uint",              // UInt64
    ]
    metrics = [
      {
        name = "rtb_wins.counter"
        type = "counter"
        tags {
          network = "{{net}}"
          source  = "{{sid}}"
        }
      }
    ]
  }

  user_info {
    store   = "statistic"
    source  = "user_info"
    target  = "stats.user_info_local"
    fields  = [
      "timemark=tm:unixnano",               // DateTime
      "aucid=auc:uuid",                     // FixedString(16)  -- Internal Auction ID
      // User IDENTITY
      "udid=udi",                           // String
      "uuid=uui:uuid",                      // FixedString(16)
      "sessid=ses:uuid",                    // FixedString(16)
      // User personal information
      "age:uint8",                          // UInt8
      "gender:uint8",                       // UInt8
      "search_gender:uint8",                // UInt8
      "email",                              // String
      "phone",                              // String
      "messanger_type",                     // String
      "messanger",                          // String
      "zip",                                // String
      "facebook",                           // String
      "twitter",                            // String
      "linkedin",                           // String
      // Targeting
      "carrier=car:uint",                   // UInt64
      "country=cc:fix*2",                   // FixedString(2)
      "city=ct:fix*5",                      // FixedString(5)
      "latitude=lt:float",                  // Float64
      "longitude=lg:float",                 // Float64
      "language=lng:fix*5",                 // FixedString(5)
    ]
    metrics = [
      {
        name = "user.counter"
        type = "counter"
        tags {
          action   = "{{e}}"
          language = "{{lng}}"
        }
      }
    ]
  }
}