import Cookies from 'js-cookie'
import {
  generateCodeChallenge,
  generateCodeVerifier,
  generateState
} from './oauth-pkce'

interface OAuthFlowOptions {
  returnTo?: string
}

const oauthCookieOptions = () => ({
  expires: 10 / 1440,
  sameSite: 'lax' as const,
  secure: location.protocol === 'https:',
  path: '/'
})

const buildAuthorizeUrl = async (
  opts: OAuthFlowOptions = {}
): Promise<string> => {
  const config = useRuntimeConfig()
  const codeVerifier = generateCodeVerifier()
  const codeChallenge = await generateCodeChallenge(codeVerifier)
  const state = generateState()
  const cookie = oauthCookieOptions()
  Cookies.set('oauth_code_verifier', codeVerifier, cookie)
  Cookies.set('oauth_state', state, cookie)
  if (opts.returnTo) {
    Cookies.set('oauth_return_to', opts.returnTo, cookie)
  }
  const params = new URLSearchParams({
    client_id: String(config.public.oidcClientId),
    redirect_uri: String(config.public.oidcRedirectUri),
    response_type: 'code',
    scope: 'openid profile email',
    state,
    code_challenge: codeChallenge,
    code_challenge_method: 'S256'
  })
  return `${config.public.oidcServerUrl}/oauth/authorize?${params}`
}

export const startOAuthLogin = async (returnTo?: string): Promise<void> => {
  window.location.href = await buildAuthorizeUrl({ returnTo })
}

export const startOAuthRegister = async (returnTo?: string): Promise<void> => {
  const config = useRuntimeConfig()
  const authorizeUrl = await buildAuthorizeUrl({ returnTo })
  const registerUrl = `${config.public.oidcFrontendUrl}/auth/register?redirect=${encodeURIComponent(authorizeUrl)}`
  window.location.href = registerUrl
}

export const consumeOAuthReturnTo = (): string | null => {
  const value = Cookies.get('oauth_return_to')
  if (value) Cookies.remove('oauth_return_to', { path: '/' })
  if (!value) return null
  try {
    const url = new URL(value, window.location.origin)
    if (url.origin === window.location.origin) {
      return url.pathname + url.search + url.hash
    }
  } catch {
    return null
  }
  return null
}
