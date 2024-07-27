export const scrollToViewCenter = (el: HTMLElement) => {
  const scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
  const { top, height } = el.getBoundingClientRect();
  // 元素的中心高度
  const elCenter = top + height / 2;
  // 窗口的中心高度
  const center = window.innerHeight / 2;
  window.scrollTo({
    top: scrollTop - (center - elCenter),
    behavior: 'smooth'
  });
}