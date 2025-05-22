<script setup>
import { ref } from 'vue'
import { t } from './i18n.js'
import ColorPicker from './components/ColorPicker.vue'
import LogoPicker from './components/LogoPicker.vue'

const text = ref('')
// Default colors: primary Razem and white background
const qrColor = ref('#720546')
const bgColor = ref('#FFFFFF')
const useLogo = ref(false)
const logoFile = ref(null)
const selectedLogo = ref(null)
const qrImageUrl = ref(null)
// Predefined swatches: Razem branding plus basic black, white and grays
const predefinedColors = ['#000000','#444444','#888888','#CCCCCC','#FFFFFF','#720546','#870f57','#aa086c','#0070CC']



async function prepareLogoFile() {
  if (selectedLogo.value) {
    const resp = await fetch(selectedLogo.value.path)
    const blob = await resp.blob()
    const name = selectedLogo.value.path.split('/').pop()
    return new File([blob], name, { type: blob.type })
  }
  return logoFile.value
}
async function generateQr() {
  if (!text.value) return
  let response
  const payload = { text: text.value, qr_color: qrColor.value, bg_color: bgColor.value }
  const logoToSend = await prepareLogoFile()
  if (logoToSend) {
    const formData = new FormData()
    formData.append('payload', JSON.stringify(payload))
    formData.append('svg_logo', logoToSend)
    response = await fetch('/api/generate-qr', { method: 'POST', body: formData })
  } else {
    response = await fetch('/api/generate-qr', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })
  }
  if (!response.ok) {
    alert(t('failed'))
    return
  }
  qrImageUrl.value = URL.createObjectURL(await response.blob())
}
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-[#720546] to-[#aa086c] flex items-center justify-center p-4">
    <div class="w-full max-w-lg bg-white bg-opacity-90 backdrop-blur-md rounded-2xl shadow-2xl p-8">
      <h1 class="text-4xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-[#720546] to-[#aa086c] mb-8 text-center">
        {{ t('title') }}
      </h1>
      <form @submit.prevent="generateQr" class="space-y-6">
        <div>
          <label for="text" class="block text-sm font-medium text-gray-700">{{ t('textOrUrl') }}</label>
          <input
            id="text"
            v-model="text"
            type="text"
            required
            :placeholder="t('enterTextUrl')"
            class="mt-1 block w-full border border-gray-300 rounded-md p-2 focus:border-[#720546] focus:outline-none"
          />
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
          <ColorPicker v-model="qrColor" :swatches="predefinedColors" :label="t('qrColorLabel')" />
          <ColorPicker v-model="bgColor" :swatches="predefinedColors" :label="t('bgColorLabel')" />
        </div>
        <div>
          <label class="inline-flex relative items-center cursor-pointer">
            <input type="checkbox" v-model="useLogo" class="sr-only peer" />
            <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-[#720546] rounded-full peer peer-checked:bg-[#720546] transition-colors"></div>
            <div class="absolute left-1 top-1 w-4 h-4 bg-white rounded-full transition-transform peer-checked:translate-x-5"></div>
            <span class="ml-3 text-gray-700">{{ t('addLogo') }}</span>
          </label>
        </div>
        <div v-if="useLogo">
          <LogoPicker v-model:selectedLogo="selectedLogo" v-model:logoFile="logoFile" />
        </div>
        <div>
          <button
            type="submit"
            class="w-full py-3 bg-gradient-to-r from-[#720546] to-[#870f57] hover:from-[#aa086c] hover:to-[#720546] text-white font-semibold rounded-lg shadow-md transition-all"
            >
            {{ t('generateQr') }}
          </button>
        </div>
      </form>
      <div v-if="qrImageUrl" class="mt-8 text-center">
        <h2 class="text-2xl font-bold mb-4 text-gray-800">{{ t('preview') }}</h2>
        <div class="inline-block p-4 bg-white rounded-lg shadow-lg">
          <img :src="qrImageUrl" alt="Generated QR code" class="block w-64 h-64" />
        </div>
        <div class="mt-4">
          <a
            :href="qrImageUrl"
            download="qrcode.png"
            class="inline-block py-2 px-6 bg-gradient-to-r from-[#0070CC] to-[#720546] hover:from-[#005fa3] hover:to-[#870f57] text-white font-semibold rounded-lg shadow-md transition-all mt-2"
            >
            {{ t('download') }}
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

