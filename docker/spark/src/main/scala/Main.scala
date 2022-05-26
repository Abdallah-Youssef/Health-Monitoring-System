package job
import org.apache.log4j.Level
import org.apache.spark.sql.SparkSession
import org.apache.spark.sql.functions._
import org.apache.spark.sql.types.{DoubleType, IntegerType, LongType, StructField, StructType, TimestampType}

object Main {
  def main(args: Array[String]): Unit = {
    org.apache.log4j.Logger.getLogger("org").setLevel(Level.FATAL)
    org.apache.log4j.Logger.getLogger("akka").setLevel(Level.FATAL)
    val spark = SparkSession.builder.master("local[*]").appName("Spark Job").getOrCreate
    import spark.implicits._


    val customSchema = StructType(Array(
      StructField("serviceName", IntegerType),
      StructField("timestamp", LongType),
      StructField("CPU", DoubleType),
      StructField("RAMTotal", DoubleType),
      StructField("RAMFree", DoubleType),
      StructField("DiskTotal", DoubleType),
      StructField("DiskFree", DoubleType)
    ))

    val logFile = spark.read
      .format("csv")
      .option("delimiter",",")
      .schema(customSchema)
      .load(args(0))

    logFile.show(2, false)

    var cal = logFile.groupBy("timestamp", "serviceName")
      .agg(
        sum("RamTotal").as("RamTotal"),
        sum("CPU").as("CPU"),
        sum("RamFree").as("RamFree"),
        sum("DiskTotal").as("DiskTotal"),
        sum("DiskFree").as("DiskFree"),
        count("serviceName").as("count"),
        max("RamTotal").as("peakRamTotal"),
        max("RamFree").as("peakRamFree"),
        max("DiskTotal").as("peakDiskTotal"),
        max("DiskFree").as("peakDiskFree")
      )

    if (args.length == 3) {
      val oldView = spark.read.parquet(args(2))
      cal = cal.union(oldView)
    }

    val newView = cal
      .groupBy("timestamp", "serviceName")
      .agg(
        sum("RamTotal").as("RamTotal"),
        sum("CPU").as("CPU"),
        sum("RamFree").as("RamFree"),
        sum("DiskTotal").as("DiskTotal"),
        sum("DiskFree").as("DiskFree"),
        sum("count").as("count"),
        max("peakRamTotal").as("peakRamTotal"),
        max("peakRamFree").as("peakRamFree"),
        max("peakDiskTotal").as("peakDiskTotal"),
        max("peakDiskFree").as("peakDiskFree")

      )
      .sort("timestamp", "serviceName")

    newView.write.format("parquet").mode("overwrite").save(args(1))

  }
}