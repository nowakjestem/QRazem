<script setup>
import { ref, computed, onMounted } from 'vue'
import { t } from './i18n.js'

const text = ref('')
const qrColor = ref('#720546')
const bgColor = ref('#0070CC')
const useLogo = ref(false)
const logoFile = ref(null)
const selectedLogo = ref(null)
const qrImageUrl = ref(null)
const showQrPicker = ref(false)
const showBgPicker = ref(false)
// Predefined swatches: Razem branding plus basic black, white and grays
const predefinedColors = [
  '#000000', '#444444', '#888888', '#CCCCCC', '#FFFFFF',
  '#720546', '#870f57', '#aa086c', '#0070CC',
]
// Predefined logos loaded from config
const predefinedLogos = ref([])
const logoQuery = ref('')

onMounted(async () => {
  try {
    const res = await fetch('/logos.json')
    predefinedLogos.value = await res.json()
  } catch (e) {
    console.error('Failed to load logos.json', e)
  }
})
const filteredLogos = computed(() => {
  const q = logoQuery.value.toLowerCase()
  return predefinedLogos.value.filter(l => l.name.toLowerCase().includes(q))
})
function selectLogo(logo) {
  selectedLogo.value = logo
  logoFile.value = null
}

function onLogoChange(event) {
  selectedLogo.value = null
  const files = event.target.files
  logoFile.value = files && files.length > 0 ? files[0] : null
}

function onLogoDrop(event) {
  selectedLogo.value = null
  const files = event.dataTransfer.files
  logoFile.value = files && files.length > 0 ? files[0] : null
}


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
            class="mt-1 block w-full bg-transparent border-b-2 border-gray-300 py-2 px-1 text-gray-800 placeholder-gray-400 focus:border-[#720546] focus:outline-none"
          />
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
          <!-- QR Color Picker -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('qrColorLabel') }}</label>
            <div class="flex items-center space-x-3">
              <div class="w-10 h-10 rounded-full border-2 border-gray-400" :style="{ backgroundColor: qrColor }"></div>
              <button
                type="button"
                @click="showQrPicker = !showQrPicker"
                class="px-3 py-1 bg-[#720546] text-white rounded-md"
              >
                {{ t('chooseColor') }}
              </button>
            </div>
            <div v-if="showQrPicker" class="mt-2 relative">
              <div class="absolute z-10 bg-white p-3 rounded-lg shadow-lg">
                <div class="grid grid-cols-4 gap-2">
                  <button v-for="c in predefinedColors" :key="c"
                    @click="qrColor = c; showQrPicker = false"
                    :style="{ backgroundColor: c }"
                    class="w-8 h-8 rounded-full border-2 border-gray-300 focus:outline-none"
                  ></button>
                </div>
                <input type="color" v-model="qrColor"
                  class="mt-2 w-full h-8 p-0 border-2 border-gray-300 rounded cursor-pointer"
                />
              </div>
            </div>
          </div>
          <!-- Background Color Picker -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('bgColorLabel') }}</label>
            <div class="flex items-center space-x-3">
              <div class="w-10 h-10 rounded-full border-2 border-gray-400" :style="{ backgroundColor: bgColor }"></div>
              <button
                type="button"
                @click="showBgPicker = !showBgPicker"
                class="px-3 py-1 bg-[#0070CC] text-white rounded-md"
              >
                {{ t('chooseColor') }}
              </button>
            </div>
            <div v-if="showBgPicker" class="mt-2 relative">
              <div class="absolute z-10 bg-white p-3 rounded-lg shadow-lg">
                <div class="grid grid-cols-4 gap-2">
                  <button v-for="c in predefinedColors" :key="c"
                    @click="bgColor = c; showBgPicker = false"
                    :style="{ backgroundColor: c }"
                    class="w-8 h-8 rounded-full border-2 border-gray-300 focus:outline-none"
                  ></button>
                </div>
                <input type="color" v-model="bgColor"
                  class="mt-2 w-full h-8 p-0 border-2 border-gray-300 rounded cursor-pointer"
                />
              </div>
            </div>
          </div>
        </div>
        <div>
          <label class="inline-flex items-center">
            <input type="checkbox" v-model="useLogo" class="form-checkbox h-5 w-5 text-[#720546]" />
            <span class="ml-2 text-gray-700">{{ t('addLogo') }}</span>
          </label>
        </div>
        <div v-if="useLogo" class="space-y-4">
          <label class="block text-sm font-medium text-gray-700">{{ t('addLogo') }}</label>
          <!-- Predefined Logos Search and Select -->
          <input
            type="text"
            v-model="logoQuery"
            placeholder="Search logos..."
            class="w-full border border-gray-300 rounded-md px-2 py-1 focus:border-[#720546] focus:outline-none"
          />
          <div class="grid grid-cols-3 sm:grid-cols-4 gap-4 max-h-48 overflow-auto">
            <div
              v-for="logo in filteredLogos"
              :key="logo.path"
              @click="selectLogo(logo)"
              :class="{'ring-2 ring-[#0070CC]': selectedLogo && selectedLogo.path === logo.path}"
              class="cursor-pointer p-1 bg-white rounded-md flex flex-col items-center"
            >
              <img :src="logo.path" :alt="logo.name" class="w-full h-16 object-contain" />
              <span class="text-sm mt-1">{{ logo.name }}</span>
            </div>
          </div>
          <!-- Or Upload Custom SVG -->
          <p class="text-sm text-gray-500">{{ t('orDragDrop') }}</p>
          <div
            class="mt-1 flex justify-center items-center px-6 pt-5 pb-6 border-2 border-dashed rounded-md cursor-pointer border-[#720546]"
            @drop.prevent="onLogoDrop"
            @dragover.prevent
          >
            <div class="space-y-2 text-center">
              <svg class="mx-auto h-12 w-12 text-[#720546]" stroke="currentColor" fill="none" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16v-4m0 0l5-5m-5 5l5 5m13-1v1a1 1 0 0 1-1 1h-5 m6-2V5a1 1 0 0 0-1-1H4a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h5m6-3h.01" />
              </svg>
              <label class="relative cursor-pointer bg-white rounded-md font-medium text-[#720546] hover:text-[#aa086c]">
                <span>{{ t('uploadSvg') }}</span>
                <input type="file" accept=".svg" class="sr-only" @change="onLogoChange" />
              </label>
              <p class="text-xs text-gray-500">{{ t('svgOnly') }}</p>
              <div v-if="logoFile" class="mt-2 text-sm text-gray-700">{{ logoFile.name }}</div>
              <div v-else-if="selectedLogo" class="mt-2 text-sm text-gray-700">{{ selectedLogo.name }}</div>
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
            class="inline-block py-2 px-6 bg-gradient-to-r from-[#0070CC] to-[#720546] hover:from-[#005fa3] hover:to-[#870f57] text-white font-semibold rounded-lg shadow-md transition-all mt-2"
            >
            {{ t('download') }}
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

