# 弹幕播放器说明

:::tip 前言
在上一代弹幕网站([Leaf](https://github.com/wangzmgit/leaf))中，
我们专门开发了弹幕播放器[WPlayer](https://github.com/wangzmgit/wplayer)，
但该播放器存在一些问题，例如仅支持vue3项目，无法直接使用video的api等。
因此，我们在这一代弹幕网站中开发了[WPlayerNext](https://github.com/wangzmgit/wplayer-next)播放器。
WPlayerNext基于开源播放器[DPlayer](https://github.com/DIYgod/DPlayer)进行了二次开发，以解决上一版本无法解决的问题。

**您也可以通过以下代码在自己的项目中引入WPlayerNext播放器，
详细功能请参考[WPlayerNext文档](https://wplayer-next.interastral-peace.com/)**
:::

### 安装

使用 npm:
```
npm install wplayer-next --save
```

使用 Yarn:
```
yarn add wplayer-next
```

### 快速开始

加载播放器文件:

```html
<div id="wplayer"></div>
<script src="WPlayer.min.js"></script>
```

或者使用模块管理器:

```js
import WPlayer from 'wplayer-next';

const player = new WPlayer(options);
```

在 js 里初始化:

```js
const player = new WPlayer({
    container: document.getElementById('wplayer'),
    video: {
        url: 'demo.mp4',
    },
});
```


