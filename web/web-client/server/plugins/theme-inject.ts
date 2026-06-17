export default defineNitroPlugin((nitro) => {
  nitro.hooks.hook('render:html', (html) => {
    html.head.unshift(
      '<script>try{var k="ui-theme-mode",m=(localStorage.getItem(k)||"light"),r=document.documentElement;r.setAttribute("data-theme",m);if(m==="dark")r.classList.add("dark");else r.classList.remove("dark")}catch(e){}</script>'
    )
  })
})
