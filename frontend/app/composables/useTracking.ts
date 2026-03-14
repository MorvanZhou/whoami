export const useTracking = () => {
  const { gtag } = useGtag()

  const trackSignUp = (method: string) => {
    gtag('event', 'sign_up', { method })
  }

  const trackLogin = (method: string) => {
    gtag('event', 'login', { method })
  }

  const trackAgentCreated = (label: string) => {
    gtag('event', 'agent_created', { agent_label: label })
  }

  const trackProfileSaved = () => {
    gtag('event', 'profile_saved')
  }

  const trackCopyPrompt = (section: string) => {
    gtag('event', 'copy_prompt', { section })
  }

  const trackCtaClick = (location: string) => {
    gtag('event', 'cta_click', { location })
  }

  return {
    trackSignUp,
    trackLogin,
    trackAgentCreated,
    trackProfileSaved,
    trackCopyPrompt,
    trackCtaClick,
  }
}
