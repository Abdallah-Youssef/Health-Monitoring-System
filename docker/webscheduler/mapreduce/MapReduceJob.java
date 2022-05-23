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

import org.apache.avro.Schema;
import org.apache.avro.SchemaBuilder;
import org.apache.avro.generic.GenericData;
import org.apache.avro.generic.GenericRecord;

import org.apache.parquet.avro.AvroParquetOutputFormat;
import org.apache.parquet.example.data.Group;


public class MapReduceJob {
    // static final Schema schema = SchemaBuilder.record("HealthMessagesSummary")
    //                                   .fields()
    //                                   .name("cpu").type().floatType().noDefault()
    //                                   .endRecord();
    static final Schema SCHEMA = new Schema.Parser().parse(
      "{\n"+
      "   \"type\": \"record\",\n"+
      "   \"name\": \"HealthMessagesSummary\",\n"+
      "   \"fields\": [\n"+
      "     {\"name\": \"cpu\", \"type\": \"double\"}\n"+
      "   ]\n" + 
      " }"
    );


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

    public static class TextToAvroParquetMapper extends Mapper<Text, HealthMessage, Void, GenericRecord> {
      
      
      GenericRecord genericRecord = new GenericData.Record(SCHEMA);

        public void map(Text serviceKey, HealthMessage healthMessage, Context context) throws IOException, InterruptedException {
              // Parse the value yourself here,
              // calling "put" on the Avro GenericRecord,
              // once for each field.  The GenericRecord
              // object is reused for every map call.
              genericRecord.put("cpu", healthMessage.cpu.get());
              context.write(null, genericRecord);
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
        job2.setMapperClass(TextToAvroParquetMapper.class);
        job2.setNumReduceTasks(0);
        job2.setOutputKeyClass(Void.class);
        job2.setOutputValueClass(Group.class);
        job2.setInputFormatClass(TextInputFormat.class);
        job2.setOutputFormatClass(AvroParquetOutputFormat.class);
        AvroParquetOutputFormat.setSchema(job2, SCHEMA);
        FileInputFormat.addInputPath(job2, new Path(args[2]));
        FileOutputFormat.setOutputPath(job2, new Path(args[3]));

        if (!job2.waitForCompletion(true)) {
            System.exit(1);
        }
    }
}


