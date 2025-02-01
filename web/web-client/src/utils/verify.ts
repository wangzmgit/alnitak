export const isEmail = (str: string) => {
  const reg = /^([a-zA-Z]|[0-9])(\w|\-)+@[a-zA-Z0-9]+\.([a-zA-Z]{2,4})$/;
  return reg.test(str);
}

export const isPhone = (str: string) => {
  const reg = /^1[3|4|5|6|7|8|9][0-9]{9}$/;
  return reg.test(str);
}

export const isUrl = (str: string) => {
  const reg = /^(http|https):\/\/[^\s/$.?#].[^\s]*$/;
  return reg.test(str);
}

export const isLegalTag = (str: string) => {
  const reg = /^[0-9a-zA-Z\u4e00-\u9fa5\uac00-\ud7ff\u0800-\u4e00]+$/
  return reg.test(str);
}