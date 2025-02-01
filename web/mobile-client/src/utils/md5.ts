import SparkMD5 from 'spark-md5';

export const getFileMD5 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const fileReader = new FileReader();
    const spark = new SparkMD5.ArrayBuffer();
    fileReader.readAsArrayBuffer(file);

    fileReader.onload = (e) => {
      if (e.target && e.target.result instanceof ArrayBuffer) {
        spark.append(e.target.result);
        resolve(spark.end());
      } else {
        reject(new Error('Failed to read file as ArrayBuffer'));
      }
    };

    fileReader.onerror = () => {
      reject(new Error('File reading failed'));
    };
  });
};