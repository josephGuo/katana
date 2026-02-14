package captcha

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdentify(t *testing.T) {
	browser := setupBrowser(t)

	tests := []struct {
		name         string
		html         string
		wantProvider Provider
		wantSiteKey  string
		wantNil      bool
	}{
		{
			name: "recaptcha v2",
			html: `<html><body>
				<div class="g-recaptcha" data-sitekey="6LcXrecapv2"></div>
				<script src="https://www.google.com/recaptcha/api.js" async defer></script>
			</body></html>`,
			wantProvider: ProviderRecaptchaV2,
			wantSiteKey:  "6LcXrecapv2",
		},
		{
			name: "recaptcha v3",
			html: `<html><body>
				<script src="https://www.google.com/recaptcha/api.js?render=6LcXrecapv3"></script>
			</body></html>`,
			wantProvider: ProviderRecaptchaV3,
			wantSiteKey:  "6LcXrecapv3",
		},
		{
			name: "cloudflare turnstile",
			html: `<html><body>
				<div class="cf-turnstile" data-sitekey="0x4AAATURNSTILE"></div>
				<script src="https://challenges.cloudflare.com/turnstile/v0/api.js" async defer></script>
			</body></html>`,
			wantProvider: ProviderTurnstile,
			wantSiteKey:  "0x4AAATURNSTILE",
		},
		{
			name: "hcaptcha",
			html: `<html><body>
				<div class="h-captcha" data-sitekey="hcap-sitekey-123"></div>
				<script src="https://js.hcaptcha.com/1/api.js" async defer></script>
			</body></html>`,
			wantProvider: ProviderHCaptcha,
			wantSiteKey:  "hcap-sitekey-123",
		},
		{
			name: "no captcha",
			html: `<html><body>
				<h1>Hello World</h1>
				<form><input type="text" name="q"><button>Search</button></form>
			</body></html>`,
			wantNil: true,
		},
		{
			name: "recaptcha v2 enterprise",
			html: `<html><body>
				<div class="g-recaptcha" data-sitekey="6LcEntV2"></div>
				<script src="https://www.google.com/recaptcha/enterprise.js" async defer></script>
			</body></html>`,
			wantProvider: ProviderRecaptchaV2Enterprise,
			wantSiteKey:  "6LcEntV2",
		},
		{
			name: "recaptcha v3 enterprise",
			html: `<html><body>
				<script src="https://www.google.com/recaptcha/enterprise.js?render=6LcEntV3"></script>
			</body></html>`,
			wantProvider: ProviderRecaptchaV3Enterprise,
			wantSiteKey:  "6LcEntV3",
		},
		{
			name: "generic data-sitekey fallback",
			html: `<html><body>
				<div id="captcha-widget" data-sitekey="generic-key-456"></div>
			</body></html>`,
			wantProvider: ProviderRecaptchaV2,
			wantSiteKey:  "generic-key-456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := servePage(t, tt.html)
			page := browser.MustPage(url)
			defer page.MustClose()
			page.MustWaitLoad()

			info, err := Identify(page)
			require.NoError(t, err)

			if tt.wantNil {
				assert.Nil(t, info)
				return
			}

			require.NotNil(t, info)
			assert.Equal(t, tt.wantProvider, info.Provider)
			assert.Equal(t, tt.wantSiteKey, info.SiteKey)
			assert.Contains(t, info.PageURL, url)
		})
	}
}
