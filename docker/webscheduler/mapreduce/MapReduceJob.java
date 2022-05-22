import java.io.IOException;
import org.json.JSONException;
import java.util.*;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.*;
import org.apache.hadoop.mapred.*;

public class MapReduceJob {

    public static class Map extends MapReduceBase implements Mapper<LongWritable, Text, Text, HealthMessage> {
        private Text serviceNameKey = new Text();

        @Override
        public void map(LongWritable _init, Text value, OutputCollector<Text, HealthMessage> outputCollector, Reporter reporter) throws IOException {
            String line = value.toString();
            String serviceName = HealthMessage.getServiceName(line);
            serviceNameKey.set(serviceName);
            try {
              HealthMessage message = new HealthMessage(line);
              outputCollector.collect(serviceNameKey, message);
            }catch (JSONException j) {
              outputCollector.collect(serviceNameKey, new HealthMessage());
            }
        }
    }

    public static class Reduce extends MapReduceBase implements Reducer<Text, HealthMessage, Text, HealthMessage> {
        @Override
        public void reduce(Text serviceKey, Iterator<HealthMessage> messages, OutputCollector<Text, HealthMessage> outputCollector, Reporter reporter) throws IOException {
            HealthMessage sum = new HealthMessage();
            while (messages.hasNext())
                sum.add(messages.next());
            outputCollector.collect(serviceKey, sum);
        }
    }


    public static void main(String[] args) throws Exception {
        JobConf conf = new JobConf(MapReduceJob.class);
        conf.setJobName("HealthMessagesJob");

        conf.setMapOutputKeyClass(Text.class);
        conf.setMapOutputValueClass(HealthMessage.class);
        conf.setOutputKeyClass(Text.class);
        conf.setOutputValueClass(HealthMessage.class);

        conf.setMapperClass(Map.class);
        conf.setCombinerClass(Reduce.class);
        conf.setReducerClass(Reduce.class);

        conf.setInputFormat(TextInputFormat.class);
        conf.setOutputFormat(TextOutputFormat.class);

        FileInputFormat.setInputPaths(conf, new Path(args[0]));
        FileOutputFormat.setOutputPath(conf, new Path(args[1]));
        JobClient.runJob(conf);
    }
}
