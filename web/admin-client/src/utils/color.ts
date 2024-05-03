export const getMainColor = (src: string) => {
  const img = document.createElement("img");
  img.src = src;
  img.crossOrigin = "anonymous";

  return new Promise((resolve, reject) => {
    img.onload = () => {
      const canvas = document.createElement("canvas");
      const ctx = canvas.getContext && canvas.getContext('2d');

      if (ctx) {
        ctx.drawImage(img, 0, 0);
        const data = ctx.getImageData(0, 0, img.width, img.height).data;
        let r = 0, g = 0, b = 0;

        //计算平均值
        for (let row = 0; row < img.height; row++) {
          for (let col = 0; col < img.width; col++) {
            r += data[(img.width * row + col) * 4];
            g += data[(img.width * row + col) * 4 + 1];
            b += data[(img.width * row + col) * 4 + 2];
          }
        }

        //计算平均值
        r /= img.width * img.height;
        g /= img.width * img.height;
        b /= img.width * img.height;

        r = Math.round(r);
        g = Math.round(g);
        b = Math.round(b);
        console.log(r, g, b)
        if (r === 0 && g === 0 && b === 0) {
          resolve("#9499a0");
        }
        resolve(`#${r.toString(16)}${g.toString(16)}${b.toString(16)}`);
      }

      reject("#9499a0");
    }
  })
} 