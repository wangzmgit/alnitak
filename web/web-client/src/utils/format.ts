import moment from "moment";

export const isEmail = (str: string) => {
    const reg = /^([a-zA-Z]|[0-9])(\w|\-)+@[a-zA-Z0-9]+\.([a-zA-Z]{2,4})$/;
    return reg.test(str);
}

export const isPhone = (str: string) => {
    const reg = /^1[3|4|5|6|7|8|9][0-9]{9}$/;
    return reg.test(str);
}

export const isUrl = (str: string) => {
    var reg = /^(http|https):\/\/[^\s/$.?#].[^\s]*$/;
    return reg.test(str);
}

export const formatDate = (time: string) => {
    return moment(time).format("YYYY-MM-DD")
}

export const formatTime = (time: string) => {
    return moment(time).format("YYYY-MM-DD hh:mm:ss")
}

export const formatRelativeTime = (timeStr: string) => {
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
        return formatTime(timeStr);
    }
}