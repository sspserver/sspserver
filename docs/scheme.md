# Service scheme

```sh
  [Advertisement Storage] := ADS (Central advertisement component)
  [Balance manager] := BM

  ADS <-> BM

  * Programm
    * DSP <- ADS
      * BidRequest handler
      * Win handler
    * SSP <- [ADS, Sourcess...]
      * SSP AdRequest hanler
      * SSP DirectAdRequest hanler
    * API methods: Click, Strict (impression, direct), Lead
```
