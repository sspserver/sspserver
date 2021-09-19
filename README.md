# Rotator

Main ads rotation tool.

## Macro constants

 * **${click_id}**      - Impression ID
 * **${domain}**        - Domain or bundle ID
 * **${zone}**          - Zone  ID
 * **${country_code}**  - Country Code ISO2
 * **${language}**      - Language code (en, fr, de, etc.)


# EXT

```sql
insert into adv_campaign(company_id, status, active, pricing_model, price, max_daily) values(1, 1, 1, 3, 1000000000, 50000000000000);

insert into adv_campaign(company_id, status, active, pricing_model, price, max_daily) values(1, 1, 1, 3, 1000000000, 50000000000000);

insert into adv_ad(campaign_id, width, height, status, active, ad_type, context, pricing_model, max_bid, max_daily) values(15, 0, 0, 1, 1, 1, '{"direct":"https://t.insigit.com/tds/int?tdsId=a0949net_r&tds_campaign=a0949net&utm_source=int&utm_campaign=1421421d&utm_content=${click_id}&data2={data2}&utm_sub=opnfnl"}', 3, 1000000000, 5000000000000);

insert into adv_ad(campaign_id, width, height, status, active, ad_type, context, pricing_model, max_bid, max_daily) values(16, 0, 0, 1, 1, 1, '{"direct":"http://udating.website/c/da57dc555e50572d?s1=2326&s2=34800&click_id=${click_id}"}', 3, 1000000000, 5000000000000);

insert into adv_zone(type, company_id, status, active, campaigns) values(1, 1, 1, 1, '{15,16}');
```
