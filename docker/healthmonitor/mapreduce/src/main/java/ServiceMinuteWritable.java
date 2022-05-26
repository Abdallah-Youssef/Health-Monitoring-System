import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.LongWritable;
import org.apache.hadoop.io.WritableComparable;

import java.io.DataInput;
import java.io.DataOutput;
import java.io.IOException;

public class ServiceMinuteWritable implements WritableComparable<ServiceMinuteWritable> {
    public IntWritable serviceName;
    public LongWritable timestamp;

    public ServiceMinuteWritable() {
        serviceName = new IntWritable();
        timestamp = new LongWritable();
    }

    public ServiceMinuteWritable(IntWritable serviceName, LongWritable timestamp) {
        this.serviceName = serviceName;
        this.timestamp = timestamp;
    }
    @Override
    public int compareTo(ServiceMinuteWritable other) {
        if(other == null)
            return 0;
        int val = serviceName.compareTo(other.serviceName);
        if(val != 0)
            return val;
        return timestamp.compareTo(other.timestamp);
    }

    @Override
    public void write(DataOutput dataOutput) throws IOException {
        serviceName.write(dataOutput);
        timestamp.write(dataOutput);
    }

    @Override
    public void readFields(DataInput dataInput) throws IOException {
        serviceName.readFields(dataInput);
        timestamp.readFields(dataInput);
    }

    @Override
    public String toString() {
        return serviceName.toString() + "," + timestamp.toString();
    }
}
