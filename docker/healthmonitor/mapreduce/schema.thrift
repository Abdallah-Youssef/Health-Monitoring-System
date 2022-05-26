struct ThriftHealthMessage{
  1: required i32 serviceName;
  2: required i64 timestamp;
  3: required double cpu;
  4: required double peakCpu;
  5: required double ramTotal;
  6: required double peakRamTotal;
  7: required double ramFree;
  8: required double peakRamFree;
  9: required double diskTotal;
  10: required double peakDiskTotal;
  11: required double diskFree;
  12: required double peakDiskFree;
  13: required i32 count;
}
