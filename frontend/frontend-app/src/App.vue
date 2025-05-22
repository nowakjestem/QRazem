<script setup>
import { ref } from 'vue'
import { t } from './i18n.js'

const text = ref('')
const qrColor = ref('#720546')
const bgColor = ref('#0070CC')
const useLogo = ref(false)
const logoFile = ref(null)
const qrImageUrl = ref(null)
const showQrPicker = ref(false)
const showBgPicker = ref(false)
const predefinedColors = ['#720546', '#870f57', '#aa086c', '#0070CC']

function onLogoChange(event) {
  const files = event.target.files
  logoFile.value = files && files.length > 0 ? files[0] : null
}

function onLogoDrop(event) {
  const files = event.dataTransfer.files
  logoFile.value = files && files.length > 0 ? files[0] : null
}

async function generateQr() {
  if (!text.value) return
  let response
  const payload = { text: text.value, qr_color: qrColor.value, bg_color: bgColor.value }
  if (logoFile.value) {
    const formData = new FormData()
    formData.append('payload', JSON.stringify(payload))
    formData.append('svg_logo', logoFile.value)
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
              <div class="w-10 h-10 rounded-full border" :style="{ backgroundColor: qrColor }"></div>
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
                  class="mt-2 w-full h-8 p-0 border-0 rounded cursor-pointer"
                />
              </div>
            </div>
          </div>
          <!-- Background Color Picker -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">{{ t('bgColorLabel') }}</label>
            <div class="flex items-center space-x-3">
              <div class="w-10 h-10 rounded-full border" :style="{ backgroundColor: bgColor }"></div>
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
                  class="mt-2 w-full h-8 p-0 border-0 rounded cursor-pointer"
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
        <div v-if="useLogo">
          <label class="block text-sm font-medium text-gray-700">{{ t('uploadSvg') }}</label>
          <div class="mt-1 flex justify-center items-center px-6 pt-5 pb-6 border-2 border-dashed rounded-md cursor-pointer border-[#720546]"
               @drop.prevent="onLogoDrop"
               @dragover.prevent>
            <div class="space-y-1 text-center">
              <svg class="mx-auto h-12 w-12 text-[#720546]" stroke="currentColor" fill="none" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M7 16v-4m0 0l5-5m-5 5l5 5m13-1v1a1 1 0 0 1-1 1h-5
                         m6-2V5a1 1 0 0 0-1-1H4a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h5m6-3h.01" />
              </svg>
              <div class="flex text-sm text-gray-600">
                <label class="relative cursor-pointer bg-white rounded-md font-medium text-[#720546] hover:text-[#aa086c]">
                  <span>{{ t('uploadSvg') }}</span>
                  <input type="file" accept=".svg" class="sr-only" @change="onLogoChange" />
                </label>
                <p class="pl-1 text-gray-700">{{ t('orDragDrop') }}</p>
              </div>
              <p class="text-xs text-gray-500">{{ t('svgOnly') }}</p>
              <div v-if="logoFile" class="mt-2 text-sm text-gray-700">{{ logoFile.name }}</div>
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

