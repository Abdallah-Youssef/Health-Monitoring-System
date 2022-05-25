struct ThriftHealthMessage{
  1: required string serviceName;
  2: required double cpu;
  3: required double peakCpu;
  4: required double ramTotal;
  5: required double peakRamTotal;
  6: required double ramFree;
  7: required double peakRamFree;
  8: required double diskTotal;
  9: required double peakDiskTotal;
  10: required double diskFree;
  11: required double peakDiskFree;
  12: required i32 count;
}