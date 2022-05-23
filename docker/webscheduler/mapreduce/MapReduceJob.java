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
import org.apache.hadoop.conf.Configuration;

import org.apache.parquet.example.data.simple.SimpleGroupFactory;
import org.apache.parquet.example.data.GroupFactory;
import org.apache.parquet.example.data.Group;
import org.apache.parquet.hadoop.example.ExampleOutputFormat;
import org.apache.parquet.schema.MessageType;
import org.apache.parquet.schema.MessageTypeParser;

public class MapReduceJob {
    static final MessageType SCHEMA = MessageTypeParser.parseMessageType(
        "message example {\n" +
          "required double cpu;\n" +
        "}");


    static final GroupFactory groupFactory = new SimpleGroupFactory(SCHEMA);

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

    public static class ThriftParquetMap extends Mapper {
      public void map(Text serviceKey, HealthMessage healthMessage, Context context) throws IOException, InterruptedException {
          Group group = groupFactory.newGroup()
                        .append("cpu", healthMessage.cpu.get());
          context.write(serviceKey, group);
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
        job1.setMapperClass(Map.class);
        job1.setCombinerClass(Reduce.class);
        job1.setReducerClass(Reduce.class);
        job1.setOutputKeyClass(Text.class);
        job1.setOutputValueClass(HealthMessage.class);
        job1.setInputFormatClass(TextInputFormat.class);
        job1.setOutputFormatClass(TextOutputFormat.class);
        FileInputFormat.addInputPath(job1, new Path(args[0]));
        FileOutputFormat.setOutputPath(job1, new Path(args[1]));

        if (!job1.waitForCompletion(true)) {
            System.exit(1);
        }

        Job job2 = Job.getInstance(conf, "ParquetFormattingJob");
        job2.setJarByClass(MapReduceJob.class);
        job2.setMapperClass(ThriftParquetMap.class);
        job2.setNumReduceTasks(0);
        job2.setOutputKeyClass(Void.class);
        job2.setOutputValueClass(Group.class);
        job2.setInputFormatClass(TextInputFormat.class);
 
        FileInputFormat.addInputPath(job2, new Path(args[2]));
        FileOutputFormat.setOutputPath(job2, new Path(args[3]));
        ExampleOutputFormat.setSchema(job2, SCHEMA);

        if (!job2.waitForCompletion(true)) {
            System.exit(1);
        }
    }
}


