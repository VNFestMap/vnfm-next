export const useKunDisableSeo = (title: string) => {
  useHead({
    htmlAttrs: { lang: 'zh-Hans' },
    meta: [
      { name: 'title', content: '' },
      { name: 'robots', content: 'noindex, nofollow' }
    ]
  })

  useSeoMeta({
    title,
    description: '',
    ogLocale: 'zh_CN',
    ogType: 'website'
  })
}
