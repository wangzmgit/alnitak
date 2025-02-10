import moment from "moment";

export const formatDate = (time: string | Date) => {
  return moment(time).format("YYYY-MM-DD")
}

export const formatTime = (time: string | Date) => {
  return moment(time).format("YYYY-MM-DD HH:mm:ss")
}

export const formatRelativeTime = (timeStr: string | Date) => {
  const currentTime = new Date();
  const time = new Date(timeStr);
  const diff = (currentTime.getTime() - time.getTime()) / 1000;
  if (diff < 60) {
    return '刚刚';
  } else if (diff < 60 * 60) {
    return Math.floor(diff / 60) + '分钟前';
  } else if (diff < 60 * 60 * 24) {
    return Math.floor(diff / (60 * 60)) + '小时前';
  } else {
    return formatDate(timeStr);
  }
}

// 获取相对时间对象
export const relativeDate = (time: string) => {
  const originalTime = new Date(time);
  const currentTime = new Date();

  const originalYear = originalTime.getFullYear()
  const originalMonth = originalTime.getMonth() + 1;
  const originalDay = originalTime.getDate();
  const hour = originalTime.getHours();
  const minute = originalTime.getMinutes();
  if (currentTime.getFullYear() !== originalYear || currentTime.getMonth() + 1 != originalMonth) {
    return {
      group: "earlier",
      time: `${originalMonth}-${originalDay} ${hour}:${minute}`
    }
  }

  switch (currentTime.getDate() - originalDay) {
    case 0:
      return {
        group: "today",
        time: `今天 ${hour}:${minute}`
      };
    case 1:
      return {
        group: "yesterday",
        time: `昨天 ${hour}:${minute}`
      };
    default:
      return {
        group: "earlier",
        time: `${originalMonth}-${originalDay} ${hour}:${minute}`
      }
  }
}

export const toDuration = (duration: number) => {
  let m: number = Math.floor(duration / 60);
  let s: number = Math.floor(duration % 60);
  const mm = m < 10 ? "0" + m : m;
  const ss = s < 10 ? "0" + s : s;
  return mm + ":" + ss;
}

// 移除html标签
export const removeHtml = (val: string) => {
  return val.replace(/<\/?.+?>/g, "");
}