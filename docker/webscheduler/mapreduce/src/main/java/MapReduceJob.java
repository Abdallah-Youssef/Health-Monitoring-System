import java.io.IOException;
import org.json.JSONException;

import java.util.*;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.*;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.Mapper;
import org.apache.hadoop.mapreduce.Reducer;
import org.apache.hadoop.mapreduce.lib.input.*;
import org.apache.hadoop.mapreduce.lib.output.*;
import org.apache.hadoop.mapreduce.lib.chain.ChainMapper;
import org.apache.hadoop.mapreduce.lib.chain.ChainReducer;
import org.apache.hadoop.conf.Configuration;

import org.apache.parquet.hadoop.thrift.ParquetThriftOutputFormat;
import org.apache.parquet.schema.MessageType;
import org.apache.parquet.schema.MessageTypeParser;

public class MapReduceJob {
    public static class Map extends Mapper<LongWritable, Text, Text, HealthMessage> {
        private Text serviceNameKey = new Text();

        public void map(LongWritable _init, Text value, Context context) throws IOException, InterruptedException {
            String line = value.toString();
            String serviceName = HealthMessage.getServiceName(line);
            serviceNameKey.set(serviceName);

            try {
              HealthMessage message = new HealthMessage(line);
              context.write(serviceNameKey, message);
            }catch (JSONException e) {
              context.write(serviceNameKey, new HealthMessage());
            }
            
        }
    }

    public static class ParquetMap extends Mapper<Text, HealthMessage, Void, ThriftHealthMessage> {
      public void map(Text serviceKey, HealthMessage healthMessage, Context context) throws IOException, InterruptedException {
          ThriftHealthMessage thriftHealthMessage = healthMessage.toThrift(serviceKey.toString());
          context.write(null, thriftHealthMessage);
      }
    }


    public static class Reduce extends Reducer<Text, HealthMessage, Text, HealthMessage> {
        public void reduce(Text serviceKey, Iterator<HealthMessage> messages, Context context) throws IOException, InterruptedException {
            HealthMessage sum = new HealthMessage();
            while (messages.hasNext())
                sum.add(messages.next());
            context.write(serviceKey, sum);
        }
    }

    public static void main(String[] args) throws Exception {    
          
        Configuration conf = new Configuration();
        
        Job job1 = Job.getInstance(conf, "HealthMessagesJob");
        job1.setJarByClass(MapReduceJob.class);

        Configuration mapConf = new Configuration(false);
        ChainMapper.addMapper(job1, Map.class, LongWritable.class, Text.class, Text.class, HealthMessage.class, mapConf);

        Configuration reduceConf = new Configuration(false);
        ChainReducer.setReducer(job1, Reduce.class, Text.class, HealthMessage.class, Text.class, HealthMessage.class, reduceConf);

        ChainReducer.addMapper(job1, ParquetMap.class, Text.class, HealthMessage.class, Void.class, ThriftHealthMessage.class, null);

        job1.setInputFormatClass(TextInputFormat.class);
        job1.setOutputFormatClass(ParquetThriftOutputFormat.class);
        FileInputFormat.addInputPath(job1, new Path(args[0]));
        ParquetThriftOutputFormat.setOutputPath(job1, new Path(args[1]));
        ParquetThriftOutputFormat.setThriftClass(job1, ThriftHealthMessage.class);

        // final FileSystem fs = FileSystem.get(conf);
        // if (fs.exists(outputPath)) fs.delete(outputPath, true);

        if (!job1.waitForCompletion(true)) {
            System.exit(1);
        }

        // Job job2 = Job.getInstance(conf, "ParquetFormattingJob");
        // job2.setJarByClass(MapReduceJob.class);
        // job2.setMapperClass(ThriftParquetMap.class);
        // job2.setNumReduceTasks(0);
        // job2.setOutputKeyClass(Void.class);
        // job2.setOutputValueClass(Group.class);
        // job2.setInputFormatClass(TextInputFormat.class);
 
        // FileInputFormat.addInputPath(job2, new Path(args[2]));
        // FileOutputFormat.setOutputPath(job2, new Path(args[3]));
        // ExampleOutputFormat.setSchema(job2, SCHEMA);

        // if (!job2.waitForCompletion(true)) {
        //     System.exit(1);
        // }
    }
}


