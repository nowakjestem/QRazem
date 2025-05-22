<template>
  <div>
    <label class="block text-sm font-medium text-gray-700 mb-1">{{ label }}</label>
    <div class="flex items-center space-x-3">
      <div
        class="w-10 h-10 rounded-full border-2 border-gray-400"
        :style="{ backgroundColor: internalValue }"
      ></div>
      <button
        type="button"
        @click="togglePicker"
        :class="buttonClasses"
      >
        {{ t('chooseColor') }}
      </button>
    </div>
    <div v-if="showPicker" class="mt-2 relative">
      <div class="absolute z-10 bg-white p-3 rounded-lg shadow-lg">
        <div class="grid grid-cols-4 gap-2">
          <button
            v-for="c in swatches"
            :key="c"
            @click="selectSwatch(c)"
            :style="{ backgroundColor: c }"
            class="w-8 h-8 rounded-full border-2 border-gray-300 focus:outline-none"
          ></button>
        </div>
        <input
          type="color"
          v-model="internalValue"
          class="mt-2 w-full h-8 p-0 border-0 rounded cursor-pointer"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { t } from '../i18n.js'

const props = defineProps({
  modelValue: { type: String, required: true },
  swatches: { type: Array, default: () => [] },
  label: { type: String, default: '' },
})
const emit = defineEmits(['update:modelValue'])

const showPicker = ref(false)
const internalValue = ref(props.modelValue)
const buttonClasses = 'px-3 py-1 bg-[#720546] text-white rounded-md'

watch(() => props.modelValue, val => {
  internalValue.value = val
})
watch(internalValue, val => {
  emit('update:modelValue', val)
})

function togglePicker() {
  showPicker.value = !showPicker.value
}

function selectSwatch(color) {
  internalValue.value = color
  showPicker.value = false
}
</script>