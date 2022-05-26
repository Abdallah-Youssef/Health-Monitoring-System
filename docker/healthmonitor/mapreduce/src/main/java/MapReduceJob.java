import java.io.IOException;
import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.*;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.FileSystem;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.*;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.Mapper;
import org.apache.hadoop.mapreduce.Reducer;
import org.apache.hadoop.mapreduce.lib.chain.ChainMapper;
import org.apache.hadoop.mapreduce.lib.chain.ChainReducer;
import org.apache.hadoop.mapreduce.lib.input.FileInputFormat;
import org.apache.hadoop.mapreduce.lib.input.TextInputFormat;
import org.apache.parquet.hadoop.thrift.ParquetThriftOutputFormat;
//import org.apache.parquet.io.*;


public class MapReduceJob {

    public static class Map extends Mapper<LongWritable, Text, ServiceMinuteWritable, HealthMessage> {
        private ServiceMinuteWritable key = new ServiceMinuteWritable();
        public void map(LongWritable _init, Text value, Context context) throws IOException, InterruptedException {
            String []tokens = value.toString().split(",");
            key.serviceName.set(Integer.parseInt(tokens[0]));

            long timeStamp = Long.parseLong(tokens[1]);
            long truncated = Instant.ofEpochSecond(timeStamp).truncatedTo(ChronoUnit.MINUTES).getEpochSecond();
            key.timestamp.set(truncated);

            HealthMessage message = new HealthMessage(tokens);
            context.write(key, message);
        }
    }

    public static class Reduce extends Reducer<ServiceMinuteWritable, HealthMessage, ServiceMinuteWritable, HealthMessage> {
        public void reduce(ServiceMinuteWritable serviceKey, Iterator<HealthMessage> messages, Context context) throws IOException, InterruptedException {
            HealthMessage sum = new HealthMessage();
            while (messages.hasNext())
                sum.add(messages.next());
            context.write(serviceKey, sum);
        }
    }

    public static class ParquetMap extends Mapper<ServiceMinuteWritable, HealthMessage, Void, ThriftHealthMessage> {
        public void map(ServiceMinuteWritable serviceKey, HealthMessage healthMessage, Context context) throws IOException, InterruptedException {
            ThriftHealthMessage thriftHealthMessage = healthMessage.toThrift(serviceKey.serviceName.get(), serviceKey.timestamp.get());
            context.write(null, thriftHealthMessage);
        }
    }



    public static void main(String[] args) throws Exception {
        String inputFiles = args[0];
        Path outputPath = new Path(args[1]);

        Configuration conf = new Configuration();

        Job job = Job.getInstance(conf, "HealthMessagesJob");
        job.setJarByClass(MapReduceJob.class);

        Configuration mapConf = new Configuration(false);
        Configuration reduceConf = new Configuration(false);

        ChainMapper.addMapper(job, Map.class, LongWritable.class, Text.class, ServiceMinuteWritable.class, HealthMessage.class, mapConf);
        ChainReducer.setReducer(job, Reduce.class, ServiceMinuteWritable.class, HealthMessage.class, ServiceMinuteWritable.class, HealthMessage.class, reduceConf);
        ChainReducer.addMapper(job, ParquetMap.class, ServiceMinuteWritable.class, HealthMessage.class, Void.class, ThriftHealthMessage.class, null);

        job.setInputFormatClass(TextInputFormat.class);
        job.setOutputFormatClass(ParquetThriftOutputFormat.class);
        FileInputFormat.addInputPaths(job, inputFiles);
        ParquetThriftOutputFormat.setOutputPath(job, outputPath);
        ParquetThriftOutputFormat.setThriftClass(job, ThriftHealthMessage.class);

        final FileSystem fs = FileSystem.get(conf);
        if (fs.exists(outputPath)) fs.delete(outputPath, true);

        if (!job.waitForCompletion(true)) {
            System.exit(1);
        }


//        JobConf conf = new JobConf(MapReduceJob.class);
//        conf.setJobName("HealthMessagesJob");
//
//        conf.setMapOutputKeyClass(ServiceMinuteWritable.class);
//        conf.setMapOutputValueClass(HealthMessage.class);
//        conf.setOutputKeyClass(ServiceMinuteWritable.class);
//        conf.setOutputValueClass(HealthMessage.class);
//
//        conf.setMapperClass(Map.class);
//        conf.setCombinerClass(Reduce.class);
//        conf.setReducerClass(Reduce.class);
//
//        conf.setInputFormat(TextInputFormat.class);
//        conf.setOutputFormat(TextOutputFormat.class);
//
//        FileInputFormat.addInputPaths(conf, args[0]);
//        FileOutputFormat.setOutputPath(conf, new Path(args[1]));
//        JobClient.runJob(conf);
    }
}
