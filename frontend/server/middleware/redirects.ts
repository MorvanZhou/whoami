export default defineEventHandler((event) => {
  const url = getRequestURL(event)
  const path = url.pathname

  if (url.hostname === 'www.whoamiagent.com') {
    return sendRedirect(event, `https://whoamiagent.com${path}${url.search}`, 301)
  }

  if (url.protocol === 'http:' && url.hostname === 'whoamiagent.com') {
    return sendRedirect(event, `https://whoamiagent.com${path}${url.search}`, 301)
  }

  if (path === '/sitemap') {
    return sendRedirect(event, '/sitemap.xml', 301)
  }

  if (path.length > 1 && path.endsWith('/') && !path.startsWith('/api/')) {
    const cleanPath = path.slice(0, -1)
    return sendRedirect(event, `${cleanPath}${url.search}`, 301)
  }
})
