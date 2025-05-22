<script setup>
import { ref, computed, watch } from 'vue'
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
// Clear logo selection or upload if toggling off
watch(useLogo, (enabled) => {
  if (!enabled) {
    selectedLogo.value = null
    logoFile.value = null
  }
})
const qrImageUrl = ref(null)
// Predefined swatches: Razem branding plus basic black, white and grays
const predefinedColors = ['#000000','#444444','#888888','#CCCCCC','#FFFFFF','#720546','#870f57','#aa086c','#0070CC']
// Download format and size options
const formats = ['svg','png','jpg']
const sizeOptions = [256, 512, 1024, 2048]
const downloadFormat = ref('png')
const downloadSize = ref(1024)
// Detect if selected logo is raster (non-SVG)
const logoIsRaster = computed(() => {
  if (selectedLogo.value) {
    return !selectedLogo.value.path.toLowerCase().endsWith('.svg')
  }
  if (logoFile.value) {
    return !logoFile.value.name.toLowerCase().endsWith('.svg')
  }
  return false
})



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
  const payload = {
    text: text.value,
    qr_color: qrColor.value,
    bg_color: bgColor.value,
    format: downloadFormat.value,
    size: downloadSize.value,
  }
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
        <div class="grid grid-cols-2 gap-4">
          <!-- Format selector -->
          <div>
            <label class="block text-sm font-medium text-gray-700">Format</label>
            <div class="relative mt-1">
              <select
                v-model="downloadFormat"
                :disabled="logoIsRaster && downloadFormat==='svg'"
                class="w-full pl-3 pr-8 py-2 border border-gray-300 bg-white rounded-md appearance-none focus:outline-none focus:border-[#720546] focus:ring-1 focus:ring-[#720546] disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <option v-for="fmt in formats" :key="fmt" :value="fmt">
                  {{ fmt.toUpperCase() }}
                </option>
              </select>
              <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2">
                <svg
                  class="h-4 w-4 text-gray-600"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </div>
            </div>
          </div>
          <!-- Size selector -->
          <div>
            <label class="block text-sm font-medium text-gray-700">Size</label>
            <div class="relative mt-1">
              <select
                v-model.number="downloadSize"
                class="w-full pl-3 pr-8 py-2 border border-gray-300 bg-white rounded-md appearance-none focus:outline-none focus:border-[#720546] focus:ring-1 focus:ring-[#720546]"
              >
                <option v-for="s in sizeOptions" :key="s" :value="s">
                  {{ s }}x{{ s }}
                </option>
              </select>
              <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2">
                <svg
                  class="h-4 w-4 text-gray-600"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </div>
            </div>
          </div>
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
            class="inline-block py-2 px-6 bg-gradient-to-r from-[#720546] to-[#aa086c] hover:from-[#870f57] hover:to-[#720546] text-white font-semibold rounded-lg shadow-md transition-all mt-2"
          >
            {{ t('download') }}
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

