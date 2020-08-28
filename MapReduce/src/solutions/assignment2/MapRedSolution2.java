package solutions.assignment2;

import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.util.StringTokenizer;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.LongWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.Mapper;
import org.apache.hadoop.mapreduce.Reducer;
import org.apache.hadoop.mapreduce.Reducer.Context;
import org.apache.hadoop.mapreduce.lib.input.FileInputFormat;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;
import org.apache.hadoop.util.GenericOptionsParser;

import examples.MapRedFileUtils;
import solutions.assignment1.MapRedSolution1.MapRecords;
import solutions.assignment1.MapRedSolution1.ReduceRecords;

public class MapRedSolution2 {

	public static class MapRecords extends Mapper<LongWritable, Text, Text, IntWritable> {
		private final static IntWritable one = new IntWritable(1);
		private Text word = new Text();

		public String iterator(String st, int j, String toIgnore, String delimitor) {
			StringTokenizer itr = new StringTokenizer(st.toString(), delimitor);
			String res = new String();
			int i = 0;
			// res=null;
			while (itr.hasMoreTokens() && (i < (j + 1))) {
				// System.out.println("i="+ i);
				String temp = itr.nextToken();
				if (i == j) {
					if (temp.compareTo(toIgnore) != 0) {
						res = temp;
					}
				}
				i++;
			}
			if (res.length() > 0) {
				return res;
			} else {
				return null;
			}
		}

		@Override
		protected void map(LongWritable key, Text value, Context context) throws IOException, InterruptedException {
			StringTokenizer itr = new StringTokenizer(value.toString(), " ");
			String res = new String();
			String res1 = new String();
			String res2 = new String();
			if (res != null) {
				// Split using ',' as its a CSV file
				res = iterator(value.toString(), 1, "tpep_pickup_datetime", ",");
				if (res != null) {
					// Next we get Date + time, separated by space and we need
					// only time,hence split using " " 2016-06-09
					// 21:06:36,2016-06-09
					res1 = iterator(res, 1, " ", " ");
					if (res != null) {
						// Next time comes like 21:06:36, hence we only need
						// first 1hour , hence split using ":"
						res2 = iterator(res1, 0, " ", ":");
						int t = Integer.parseInt(res2);
						String finalres = new String();
						// Time comes with 24 hour clock hence we need to
						// convert it into 12 hour am/pm clock
						if (t > 12) {
							t = t - 12;
							finalres = t + "pm";
						} else if (t == 0) {
							finalres = 12 + "am";
						} else {
							finalres = t + "am";
						}
						word.set(finalres);
						context.write(word, one);
					}
				}
			}
		}
	}

	public static class ReduceRecords extends Reducer<Text, IntWritable, Text, IntWritable> {
		private IntWritable result = new IntWritable();

		@Override
		protected void reduce(Text key, Iterable<IntWritable> values, Context context)
				throws IOException, InterruptedException {
			int sum = 0;

			for (IntWritable val : values)
				sum += val.get();

			result.set(sum);
			context.write(key, result);
		}
	}

	public static void main(String[] args) throws Exception {
		Configuration conf = new Configuration();

		String[] otherArgs = new GenericOptionsParser(conf, args).getRemainingArgs();

		if (otherArgs.length != 2) {
			System.err.println("Usage: MapRedSolution2 <in> <out>");
			System.exit(2);
		}

		Job job = Job.getInstance(conf, "MapRed Solution #2");

		job.setMapperClass(MapRecords.class);
		job.setCombinerClass(ReduceRecords.class);
		job.setReducerClass(ReduceRecords.class);

		job.setOutputKeyClass(Text.class);
		job.setOutputValueClass(IntWritable.class);

		FileInputFormat.addInputPath(job, new Path(otherArgs[0]));
		FileOutputFormat.setOutputPath(job, new Path(otherArgs[1]));

		MapRedFileUtils.deleteDir(otherArgs[1]);
		int exitCode = job.waitForCompletion(true) ? 0 : 1;

		FileInputStream fileInputStream = new FileInputStream(new File(otherArgs[1] + "/part-r-00000"));
		String md5 = org.apache.commons.codec.digest.DigestUtils.md5Hex(fileInputStream);
		fileInputStream.close();

		String[] validMd5Sums = { "03357cb042c12da46dd5f0217509adc8", "ad6697014eba5670f6fc79fbac73cf83",
				"07f6514a2f48cff8e12fdbc533bc0fe5", "e3c247d186e3f7d7ba5bab626a8474d7",
				"fce860313d4924130b626806fa9a3826", "cc56d08d719a1401ad2731898c6b82dd",
				"6cd1ad65c5fd8e54ed83ea59320731e9", "59737bd718c9f38be5354304f5a36466",
				"7d35ce45afd621e46840627a79f87dac" };

		for (String validMd5 : validMd5Sums) {
			if (validMd5.contentEquals(md5)) {
				System.out.println("The result looks good :-)");
				System.exit(exitCode);
			}
		}
		System.out.println("The result does not look like what we expected :-(");
		System.exit(exitCode);
	}
}
