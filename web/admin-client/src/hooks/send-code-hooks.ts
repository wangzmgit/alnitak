import { ref } from 'vue';
import { useMessage } from 'naive-ui';
import { sendEmailCodeAPI } from '@/api/code';
import { statusCode } from '@/utils/status-code';

export default function useSendCode() {

  const disabledSend = ref(false);//禁用发送按钮
  const sendBtnText = ref('发送验证码');//发送按钮文字
  const message = useMessage();//通知

  const sendEmailCodeAsync = async (email: string, captchaId: string) => {
    //禁用发送按钮
    disabledSend.value = true;
    try {
      const res = await sendEmailCodeAPI({ email, captchaId })
      if (res.data.code === statusCode.OK) {
        message.success('发送成功');
        //开启倒计时
        let count = 0;
        let tag = setInterval(() => {
          if (++count >= 60) {
            clearInterval(tag);
            disabledSend.value = false;
            sendBtnText.value = '发送验证码';
            return;
          }
          sendBtnText.value = `${60 - count}秒`;
        }, 1000);
      } else {
        disabledSend.value = false;
      }
      return res.data;
    } catch {
      disabledSend.value = false;
      sendBtnText.value = '发送验证码';
      message.error('发送失败');
    }
  }

  return {
    disabledSend,
    sendBtnText,
    sendEmailCodeAsync
  }
}

