import org.apache.hadoop.io.DoubleWritable;
import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.Writable;
import org.json.JSONObject;
import org.json.JSONException;

import java.io.DataInput;
import java.io.DataOutput;
import java.io.IOException;

public class HealthMessage implements Writable {
    public DoubleWritable cpu;
    public DoubleWritable peakCpu;
    public DoubleWritable ramTotal;
    public DoubleWritable peakRamTotal;
    public DoubleWritable ramFree;
    public DoubleWritable peakRamFree;
    public DoubleWritable diskTotal;
    public DoubleWritable peakDiskTotal;
    public DoubleWritable diskFree;
    public DoubleWritable peakDiskFree;
    public IntWritable count;

    public HealthMessage(String message) throws JSONException {
        JSONObject object = new JSONObject(message);
        cpu = new DoubleWritable(object.getDouble("CPU"));
        ramTotal = new DoubleWritable(object.getJSONObject("RAM").getDouble("Total"));
        ramFree = new DoubleWritable(object.getJSONObject("RAM").getDouble("Free"));
        diskTotal = new DoubleWritable(object.getJSONObject("Disk").getDouble("Total"));
        diskFree = new DoubleWritable(object.getJSONObject("Disk").getDouble("Free"));
        count = new IntWritable(1);
        peakCpu = new DoubleWritable(0);
        peakRamTotal = new DoubleWritable(0);
        peakRamFree = new DoubleWritable(0);
        peakDiskTotal = new DoubleWritable(0);
        peakDiskFree = new DoubleWritable(0);
    }

    public HealthMessage() {
        cpu = new DoubleWritable(0.0);
        ramTotal = new DoubleWritable(0.0);
        ramFree = new DoubleWritable(0.0);
        diskTotal = new DoubleWritable(0.0);
        diskFree = new DoubleWritable(0.0);
        count = new IntWritable(0);
        peakCpu = new DoubleWritable(0);
        peakRamTotal = new DoubleWritable(0);
        peakRamFree = new DoubleWritable(0);
        peakDiskTotal = new DoubleWritable(0);
        peakDiskFree = new DoubleWritable(0);

    }

    public void add(HealthMessage other) {
        cpu.set(cpu.get() + other.cpu.get());
        ramTotal.set(ramTotal.get() + other.ramTotal.get());
        ramFree.set(ramFree.get() + other.ramFree.get());
        diskTotal.set(diskTotal.get() + other.diskTotal.get());
        diskFree.set(diskFree.get() + other.diskFree.get());
        count.set(count.get() + other.count.get());
        peakCpu.set(Math.max(peakCpu.get(), other.cpu.get()));
        peakRamTotal.set(Math.max(peakRamTotal.get(), other.ramTotal.get()));
        peakRamFree.set(Math.max(peakRamFree.get(), other.ramFree.get()));
        peakDiskTotal.set(Math.max(peakDiskTotal.get(), other.diskTotal.get()));
        peakDiskFree.set(Math.max(peakDiskFree.get(), other.diskFree.get()));
    }


    public static String getServiceName(String message){
      try {
        JSONObject object = new JSONObject(message);
        return object.getString("serviceName");
      }
      catch (JSONException e) {
        return "JSONException";
      }
    }


    @Override
    public String toString() {
        int n = count.get();
        return  cpu.get() / n + "," +
                ramTotal.get() / n + "," +
                ramFree.get() / n + "," +
                diskTotal.get() / n + "," +
                diskFree.get() / n + "," +
                peakCpu.get() + "," +
                peakRamTotal.get() + "," +
                peakRamFree.get() + "," +
                peakDiskTotal.get() + "," +
                peakDiskFree.get();
    }


    @Override
    public void write(DataOutput dataOutput) throws IOException {
        cpu.write(dataOutput);
        ramFree.write(dataOutput);
        ramTotal.write(dataOutput);
        diskFree.write(dataOutput);
        diskTotal.write(dataOutput);
        count.write(dataOutput);
    }

    @Override
    public void readFields(DataInput dataInput) throws IOException {
        cpu.readFields(dataInput);
        ramFree.readFields(dataInput);
        ramTotal.readFields(dataInput);
        diskFree.readFields(dataInput);
        diskTotal.readFields(dataInput);
        count.readFields(dataInput);
    }
}
