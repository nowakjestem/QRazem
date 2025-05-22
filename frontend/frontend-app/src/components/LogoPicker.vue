<template>
  <div class="space-y-4">
    <label class="block text-sm font-medium text-gray-700">{{ t('addLogo') }}</label>
    <input
      type="text"
      v-model="query"
      :placeholder="t('searchLogos')"
      class="w-full border border-gray-300 rounded-md p-2 focus:border-[#720546] focus:outline-none"
    />
    <div class="grid grid-cols-3 sm:grid-cols-4 gap-4 max-h-48 overflow-auto">
      <div
        v-for="logo in filtered"
        :key="logo.path"
        @click="onSelectLogo(logo)"
        :class="[
          'cursor-pointer bg-white rounded-md p-2 flex flex-col items-center',
          isSelectedLogo(logo) ? 'border-2 border-[#0070CC]' : ''
        ]"
      >
        <div class="w-16 h-16 bg-white flex items-center justify-center">
          <img
            :src="logo.path"
            :alt="logo.name"
            class="max-w-full max-h-full object-contain"
          />
        </div>
        <span class="text-sm mt-1 block text-center">{{ logo.name }}</span>
      </div>
    </div>
    <!-- Preview selected predefined logo -->
    <div v-if="selectedLogo" class="flex items-center space-x-3 mt-2 border-2 border-[#0070CC] bg-white rounded-md p-2">
      <div class="w-12 h-12 bg-white rounded flex items-center justify-center">
        <img :src="selectedLogo.path" :alt="selectedLogo.name" class="max-w-full max-h-full object-contain" />
      </div>
      <span class="text-md font-medium text-gray-800">{{ selectedLogo.name }}</span>
    </div>
    <p class="text-sm text-gray-500">{{ t('orDragDrop') }}</p>
    <div
      class="mt-1 flex justify-center items-center px-6 pt-5 pb-6 border-2 border-dashed rounded-md cursor-pointer border-[#720546]"
      @drop.prevent="onFileDrop"
      @dragover.prevent
    >
      <div class="space-y-2 text-center">
        <!-- Upload icon -->
        <svg xmlns="http://www.w3.org/2000/svg" class="mx-auto h-12 w-12 text-[#720546]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M4 16v1a2 2 0 002 2h12a2 2 0 002-2v-1m-4-4l-4-4m0 0l-4 4m4-4v12" />
        </svg>
        <label class="cursor-pointer font-medium text-[#720546] hover:text-[#aa086c]">
          <span>{{ t('uploadSvg') }}</span>
          <input type="file" accept="image/*" class="sr-only" @change="onFileSelect" />
        </label>
        <p class="text-xs text-gray-500">{{ t('imageOnly') }}</p>
        <div v-if="logoFile" class="mt-2 text-sm text-gray-700">{{ logoFile.name }}</div>
        <div v-else-if="selectedLogo" class="mt-2 text-sm text-gray-700">{{ selectedLogo.name }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { t } from '../i18n.js'

const props = defineProps({
  modelSelectedLogo: Object,
  modelLogoFile: Object,
})
const emit = defineEmits(['update:modelSelectedLogo', 'update:modelLogoFile'])

const logos = ref([])
const query = ref('')
onMounted(async () => {
  try {
    const res = await fetch('/logos.json')
    logos.value = await res.json()
  } catch (e) {
    console.error('Failed to load logos.json', e)
  }
})
const filtered = computed(() =>
  logos.value.filter(l => l.name.toLowerCase().includes(query.value.toLowerCase()))
)
const selectedLogo = computed({
  get: () => props.modelSelectedLogo,
  set: v => emit('update:modelSelectedLogo', v),
})
const logoFile = computed({
  get: () => props.modelLogoFile,
  set: v => emit('update:modelLogoFile', v),
})

function onSelectLogo(logo) {
  selectedLogo.value = logo
  logoFile.value = null
}
function isSelectedLogo(logo) {
  return props.modelSelectedLogo && props.modelSelectedLogo.path === logo.path
}
function onFileSelect(event) {
  selectedLogo.value = null
  const files = event.target.files
  logoFile.value = files && files.length > 0 ? files[0] : null
}
function onFileDrop(event) {
  selectedLogo.value = null
  const files = event.dataTransfer.files
  logoFile.value = files && files.length > 0 ? files[0] : null
}
</script>