<script setup lang="ts">
import { ref } from 'vue';
import { useAppStore } from '@/store';
import { Search, RefreshCw, X } from 'lucide-vue-next';
import TimeRangeSelector from './TimeRangeSelector.vue';
import ServiceSelector from './ServiceSelector.vue';

const store = useAppStore();
const loading = ref(false);

// Individual filter states
const selectedService = ref('');
const selectedLevel = ref('');
const hostFilter = ref('');
const messageFilter = ref('');
const traceIdFilter = ref('');
const clientIpFilter = ref('');

const emit = defineEmits(['search', 'refresh', 'auto-refresh']);
const currentRange = ref('15m');

const hasFilters = () => {
  return selectedService.value || selectedLevel.value || hostFilter.value ||
    messageFilter.value || traceIdFilter.value || clientIpFilter.value;
};

const buildQuery = () => {
  const parts: string[] = [];

  if (selectedService.value) {
    parts.push(`service:${selectedService.value}`);
  }
  if (selectedLevel.value) {
    parts.push(`level:${selectedLevel.value}`);
  }
  if (hostFilter.value.trim()) {
    parts.push(`host:${hostFilter.value.trim()}`);
  }
  if (traceIdFilter.value.trim()) {
    parts.push(`trace_id:${traceIdFilter.value.trim()}`);
  }
  if (clientIpFilter.value.trim()) {
    parts.push(`client_ip:${clientIpFilter.value.trim()}`);
  }
  if (messageFilter.value.trim()) {
    const msg = messageFilter.value.trim();
    // Wrap in quotes if contains spaces or special chars treated as tokens by the lexer
    const needsQuotes = /[\s"():!]/.test(msg);
    if (needsQuotes) {
      // Escape backslash and double-quote for the lexer's readString() handler
      const escaped = msg.replace(/\\/g, '\\\\').replace(/"/g, '\\"');
      parts.push(`"${escaped}"`);
    } else {
      parts.push(msg);
    }
  }

  return parts.join(' AND ');
};

const doSearch = () => {
  const query = buildQuery();
  emit('search', query, currentRange.value);
};

const handleRefresh = () => {
  const query = buildQuery();
  emit('refresh', query, currentRange.value);
};

const handleRangeUpdate = (range: string) => {
  currentRange.value = range;
  const query = buildQuery();
  emit('search', query, range);
};

const handleServiceSelect = (service: string) => {
  selectedService.value = service;
  const query = buildQuery();
  emit('search', query, currentRange.value);
};

const handleLevelChange = () => {
  const query = buildQuery();
  emit('search', query, currentRange.value);
};

const handleKeyEnter = () => {
  doSearch();
};

const clearAllFilters = () => {
  selectedService.value = '';
  selectedLevel.value = '';
  hostFilter.value = '';
  messageFilter.value = '';
  traceIdFilter.value = '';
  clientIpFilter.value = '';
  doSearch();
};

const levelOptions = ['DEBUG', 'INFO', 'WARN', 'ERROR', 'FATAL'];
</script>

<template>
  <header class="min-h-16 bg-gray-900 border-b border-gray-800 flex items-center px-8 py-3 shrink-0 relative z-40 flex-wrap gap-y-2">
    <!-- Row 1: Time + Service + Level + Host + TraceID + ClientIP + Clear/Refresh -->
    <div class="flex items-center space-x-3 flex-wrap gap-y-2 w-full">
      <TimeRangeSelector @update="handleRangeUpdate" @auto-refresh="$emit('auto-refresh', $event)" />
      <ServiceSelector @select="handleServiceSelect" />

      <!-- Level Dropdown (with clear) -->
      <div class="flex items-center space-x-1">
        <select
          v-model="selectedLevel"
          @change="handleLevelChange"
          class="bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-300 focus:outline-none focus:ring-2 focus:ring-cyan-500/50 focus:border-cyan-500 transition-all cursor-pointer"
          :class="selectedLevel ? 'text-cyan-400 border-cyan-500/50' : ''"
        >
          <option value="">{{ store.t('search.level_all') }}</option>
          <option v-for="lvl in levelOptions" :key="lvl" :value="lvl">{{ lvl }}</option>
        </select>
        <button
          v-if="selectedLevel"
          type="button"
          @click="selectedLevel = ''; doSearch()"
          class="p-1 rounded-md text-gray-500 hover:text-gray-300 hover:bg-gray-700 transition-colors z-10"
          title="Clear level"
        >
          <X class="w-3 h-3" />
        </button>
      </div>

      <!-- Host -->
      <div class="relative group">
        <input
          type="text"
          v-model="hostFilter"
          @keyup.enter="handleKeyEnter"
          :placeholder="store.t('search.host_placeholder')"
          class="w-36 bg-gray-800/70 border border-gray-700/50 rounded-lg pl-3 pr-7 py-1.5 text-xs text-gray-300 placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-cyan-500/50 focus:border-cyan-500 transition-all"
        />
        <button
          v-if="hostFilter"
          type="button"
          @click="hostFilter = ''; doSearch()"
          class="absolute right-1.5 top-1/2 -translate-y-1/2 p-1 rounded text-gray-600 hover:text-gray-300 hover:bg-gray-700 transition-colors z-10"
          title="Clear"
        >
          <X class="w-3 h-3" />
        </button>
      </div>

      <!-- Trace ID -->
      <div class="relative group">
        <input
          type="text"
          v-model="traceIdFilter"
          @keyup.enter="handleKeyEnter"
          :placeholder="store.t('search.trace_id_placeholder')"
          class="w-44 bg-gray-800/70 border border-gray-700/50 rounded-lg pl-3 pr-7 py-1.5 text-xs text-gray-300 placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-cyan-500/50 focus:border-cyan-500 transition-all"
        />
        <button
          v-if="traceIdFilter"
          type="button"
          @click="traceIdFilter = ''; doSearch()"
          class="absolute right-1.5 top-1/2 -translate-y-1/2 p-1 rounded text-gray-600 hover:text-gray-300 hover:bg-gray-700 transition-colors z-10"
          title="Clear"
        >
          <X class="w-3 h-3" />
        </button>
      </div>

      <!-- Client IP -->
      <div class="relative group">
        <input
          type="text"
          v-model="clientIpFilter"
          @keyup.enter="handleKeyEnter"
          :placeholder="store.t('search.client_ip_placeholder')"
          class="w-36 bg-gray-800/70 border border-gray-700/50 rounded-lg pl-3 pr-7 py-1.5 text-xs text-gray-300 placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-cyan-500/50 focus:border-cyan-500 transition-all"
        />
        <button
          v-if="clientIpFilter"
          type="button"
          @click="clientIpFilter = ''; doSearch()"
          class="absolute right-1.5 top-1/2 -translate-y-1/2 p-1 rounded text-gray-600 hover:text-gray-300 hover:bg-gray-700 transition-colors z-10"
          title="Clear"
        >
          <X class="w-3 h-3" />
        </button>
      </div>

      <!-- Clear + Refresh -->
      <button
        v-if="hasFilters()"
        type="button"
        @click="clearAllFilters"
        class="flex items-center space-x-1 px-2 py-1 text-xs text-gray-500 hover:text-gray-300 transition-colors"
        :title="store.t('search.clear_filters')"
      >
        <X class="w-3 h-3" />
        <span>{{ store.t('search.clear_filters') }}</span>
      </button>
      <button @click="handleRefresh" class="p-2 hover:bg-gray-800 rounded-lg text-gray-400 hover:text-white transition-all">
        <RefreshCw class="w-5 h-5" :class="{'animate-spin': loading}" />
      </button>
    </div>

    <!-- Row 2: Message + Search button -->
    <div class="flex items-center space-x-3 w-full">
      <div class="relative flex-1 min-w-[200px] group">
        <input
          type="text"
          v-model="messageFilter"
          @keyup.enter="handleKeyEnter"
          :placeholder="store.t('search.message_placeholder')"
          class="w-full bg-gray-800 border border-gray-700 rounded-lg pl-3 pr-8 py-2 text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-cyan-500/50 focus:border-cyan-500 transition-all"
        />
        <button
          v-if="messageFilter"
          type="button"
          @click="messageFilter = ''; doSearch()"
          class="absolute right-2 top-1/2 -translate-y-1/2 p-0.5 rounded text-gray-500 hover:text-gray-300 hover:bg-gray-700 transition-colors z-10"
          title="Clear"
        >
          <X class="w-3.5 h-3.5" />
        </button>
      </div>

      <button
        @click="doSearch"
        class="flex items-center space-x-1.5 px-4 py-1.5 bg-cyan-600 hover:bg-cyan-500 text-white text-xs font-bold rounded-lg transition-all active:scale-95 shadow-lg"
      >
        <Search class="w-3.5 h-3.5" />
        <span>{{ store.t('search.search_btn') }}</span>
      </button>
    </div>
  </header>
</template>
