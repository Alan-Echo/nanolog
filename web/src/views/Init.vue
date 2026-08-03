<script setup lang="ts">
import { ref, computed } from 'vue';
import { useAppStore, useToastStore } from '@/store';
import { api } from '@/api';
import { Shield, Lock, User as UserIcon, Loader2, Languages, ArrowRight, Eye, EyeOff, Server } from 'lucide-vue-next';

const store = useAppStore();
const toastStore = useToastStore();
const form = ref({ username: '', password: '', confirmPassword: '' });
const showPassword = ref(false);
const showConfirmPassword = ref(false);
const loading = ref(false);
const error = ref('');

const t = computed(() => store.t);
const currentLang = computed(() => store.currentLang);

const toggleLang = () => {
    const nextLang = store.currentLang === 'zh' ? 'en' : 'zh';
    store.setLang(nextLang);
};

const handleInit = async () => {
    if (loading.value) return;

    // Client-side validation
    if (!form.value.username.trim()) {
        error.value = store.currentLang === 'zh' ? '请输入管理员用户名' : 'Please enter an admin username';
        return;
    }
    if (!form.value.password) {
        error.value = store.currentLang === 'zh' ? '请输入密码' : 'Please enter a password';
        return;
    }
    if (form.value.password.length < 6) {
        error.value = store.currentLang === 'zh' ? '密码长度至少6位' : 'Password must be at least 6 characters';
        return;
    }
    if (form.value.password !== form.value.confirmPassword) {
        error.value = store.currentLang === 'zh' ? '两次输入的密码不一致' : 'Passwords do not match';
        return;
    }

    loading.value = true;
    error.value = '';
    try {
        const data = await api.init({ username: form.value.username, password: form.value.password });
        store.setAuth(data.token, data.username, data.role, false);
        store.setSystemStatus({ node_role: store.nodeRole, version: store.systemVersion, initialized: true });
        toastStore.showToast(store.currentLang === 'zh' ? '系统初始化成功！' : 'System initialized successfully!', 'success');
    } catch (e: any) {
        const msg = e.message || (store.currentLang === 'zh' ? '初始化失败' : 'Initialization failed');
        error.value = msg;
        toastStore.showToast(msg, 'error');
    } finally {
        loading.value = false;
    }
};
</script>

<template>
  <div class="min-h-screen w-full flex flex-col lg:flex-row bg-[#020617] text-slate-200 selection:bg-emerald-500/30 overflow-x-hidden">
    <!-- Fixed Noise Texture Overlay -->
    <div class="fixed inset-0 pointer-events-none z-[100] opacity-[0.03] mix-blend-overlay bg-[url('https://grainy-gradients.vercel.app/noise.svg')]"></div>

    <!-- Language Switcher -->
    <button
        @click.stop="toggleLang"
        type="button"
        class="fixed top-6 right-6 z-[110] flex items-center space-x-2 px-5 py-2.5 bg-slate-900/60 hover:bg-slate-800/80 backdrop-blur-xl border border-slate-700/50 rounded-full transition-all duration-300 group shadow-2xl active:scale-95"
    >
        <Languages class="w-4 h-4 text-emerald-400 group-hover:rotate-12 transition-transform duration-500" />
        <span class="text-[11px] font-black uppercase tracking-[0.2em] text-slate-300 group-hover:text-white">
            {{ currentLang === 'zh' ? 'En' : '中文' }}
        </span>
    </button>

    <!-- Left Side: Branding & Welcome -->
    <div class="relative w-full lg:w-[60%] xl:w-[65%] flex flex-col justify-start p-8 lg:p-24 pt-20 lg:pt-32 border-b lg:border-b-0 lg:border-r border-slate-800/20 min-h-[50vh] lg:min-h-screen">
      <!-- Animated Mesh Gradient Background -->
      <div class="absolute inset-0 overflow-hidden pointer-events-none">
        <div class="absolute top-[-10%] left-[-5%] w-[60%] h-[60%] bg-emerald-600/10 rounded-full blur-[120px] animate-pulse"></div>
        <div class="absolute bottom-[-10%] right-[-5%] w-[50%] h-[50%] bg-teal-500/10 rounded-full blur-[100px] animate-pulse" style="animation-delay: 2s;"></div>
      </div>

      <!-- Content -->
      <div class="relative z-10 max-w-2xl animate-in fade-in slide-in-from-left-8 duration-1000">
        <div class="flex items-center space-x-5 mb-10">
            <div class="relative group cursor-pointer">
                <div class="absolute -inset-2 bg-gradient-to-tr from-emerald-500/20 to-teal-600/20 rounded-2xl blur-lg group-hover:opacity-100 transition duration-700 opacity-50"></div>
            <div class="relative w-16 h-16 flex items-center justify-center overflow-hidden rounded-2xl border border-slate-800/10">
                <img src="/logo.png" alt="NanoLog" class="w-full h-full object-contain transform group-hover:scale-110 transition-transform duration-500" />
            </div>
            </div>
            <div class="h-10 w-px bg-slate-800/50"></div>
            <span class="text-3xl font-black tracking-tighter text-white uppercase italic">Nano<span class="text-emerald-500">Log</span></span>
        </div>

        <h1 class="text-5xl lg:text-7xl font-black text-white leading-[1.1] mb-6 tracking-tight">
          <span class="relative inline-block mt-4">
            <span class="bg-clip-text text-transparent bg-gradient-to-r from-emerald-400 via-teal-500 to-cyan-500 italic whitespace-nowrap">{{ t('init.title') }}</span>
            <div class="absolute -bottom-2 left-0 w-3/4 h-1.5 bg-emerald-500/30 rounded-full blur-[2px]"></div>
          </span>
        </h1>

        <p class="text-slate-400 text-lg leading-relaxed max-w-xl font-medium">
            {{ t('init.subtitle') }}
        </p>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mt-12">
            <div class="group p-8 rounded-3xl bg-slate-900/20 border border-slate-800/30 hover:border-emerald-500/30 transition-all duration-500 backdrop-blur-sm">
                <div class="w-12 h-12 rounded-2xl bg-emerald-500/10 flex items-center justify-center mb-6 group-hover:bg-emerald-500/20 transition-colors">
                    <Shield class="w-6 h-6 text-emerald-400" />
                </div>
                <h3 class="text-xl font-bold text-white mb-3">{{ t('auth.feature_security_title') }}</h3>
                <p class="text-slate-500 text-sm leading-relaxed font-medium">{{ t('auth.feature_security_desc') }}</p>
            </div>

            <div class="group p-8 rounded-3xl bg-slate-900/20 border border-slate-800/30 hover:border-teal-500/30 transition-all duration-500 backdrop-blur-sm">
                <div class="w-12 h-12 rounded-2xl bg-teal-500/10 flex items-center justify-center mb-6 group-hover:bg-teal-500/20 transition-colors">
                    <Server class="w-6 h-6 text-teal-400" />
                </div>
                <h3 class="text-xl font-bold text-white mb-3">AES-GCM + Bcrypt</h3>
                <p class="text-slate-500 text-sm leading-relaxed font-medium">{{ store.currentLang === 'zh' ? '所有凭据均经过加密和哈希处理，确保静态数据安全。' : 'All credentials encrypted and hashed for security at rest.' }}</p>
            </div>
        </div>

        <div class="mt-12 flex items-center space-x-12 opacity-40">
            <div class="flex flex-col">
                <span class="text-3xl font-black text-white italic tracking-tighter">{{ store.systemVersion }}</span>
                <span class="text-[10px] uppercase tracking-[0.2em] font-black text-slate-500">{{ t('auth.release_status') }}</span>
            </div>
            <div class="h-10 w-px bg-slate-800"></div>
            <div class="flex flex-col">
                <span class="text-3xl font-black text-white italic tracking-tighter">AES-256</span>
                <span class="text-[10px] uppercase tracking-[0.2em] font-black text-slate-500">{{ store.currentLang === 'zh' ? '加密标准' : 'Encryption' }}</span>
            </div>
        </div>
      </div>
    </div>

    <!-- Right Side: Initialization Form -->
    <div class="relative w-full lg:w-[40%] xl:w-[35%] flex flex-col justify-start p-8 lg:p-16 pt-12 lg:pt-10 bg-[#020617] lg:bg-transparent min-h-screen">
        <!-- Animated vertical separator -->
        <div class="hidden lg:block absolute left-0 top-0 bottom-0 w-px bg-gradient-to-b from-transparent via-slate-800/50 to-transparent overflow-hidden">
            <div class="absolute top-0 left-0 w-full h-32 bg-gradient-to-b from-transparent via-emerald-500/50 to-transparent blur-[1px] animate-glow-scan"></div>
        </div>

        <div class="w-full max-w-md mx-auto animate-in fade-in slide-in-from-right-12 duration-1000">
            <div class="mb-10">
                <div class="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-[10px] font-black uppercase tracking-[0.2em] mb-8">
                    <div class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></div>
                    <span>{{ store.currentLang === 'zh' ? '首次启动' : 'First Launch' }}</span>
                </div>
                <h2 class="text-4xl lg:text-5xl font-black text-white tracking-tight mb-4 italic">{{ t('init.title') }}</h2>
                <p class="text-slate-500 font-medium text-base tracking-wide">{{ t('init.subtitle') }}</p>
            </div>

            <!-- Form Section -->
            <form @submit.prevent="handleInit" class="space-y-6">
                <Transition
                    enter-active-class="transition duration-500 cubic-bezier(0.16, 1, 0.3, 1)"
                    enter-from-class="transform -translate-y-4 opacity-0"
                    enter-to-class="transform translate-y-0 opacity-100"
                    leave-active-class="transition duration-300 ease-in"
                    leave-from-class="opacity-100"
                    leave-to-class="opacity-0"
                >
                    <div v-if="error" class="bg-red-500/10 border border-red-500/20 rounded-2xl p-5 text-xs text-red-400 font-bold flex items-center space-x-4">
                        <div class="w-2 h-2 rounded-full bg-red-500 animate-pulse"></div>
                        <span class="tracking-wide">{{ error }}</span>
                    </div>
                </Transition>

                <!-- Username -->
                <div class="space-y-4 group">
                    <label class="text-[11px] uppercase tracking-[0.3em] font-black text-slate-500 group-focus-within:text-emerald-500 transition-colors block ml-1">
                        {{ t('init.admin_username') }}
                    </label>
                    <div class="relative">
                        <input
                            type="text"
                            v-model="form.username"
                            required
                            autocomplete="username"
                            placeholder="admin"
                            class="w-full h-16 bg-slate-900/40 border border-slate-800/60 rounded-xl px-6 text-base text-white placeholder-slate-600 focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500/50 transition-all hover:bg-slate-900/60"
                        />
                        <div class="absolute right-6 top-1/2 -translate-y-1/2">
                            <UserIcon class="w-5 h-5 text-slate-600 group-focus-within:text-emerald-500 transition-colors" />
                        </div>
                    </div>
                </div>

                <!-- Password -->
                <div class="space-y-4 group">
                    <label class="text-[11px] uppercase tracking-[0.3em] font-black text-slate-500 group-focus-within:text-emerald-500 transition-colors block ml-1">
                        {{ t('init.root_password') }}
                    </label>
                    <div class="relative">
                        <input
                            :type="showPassword ? 'text' : 'password'"
                            v-model="form.password"
                            required
                            autocomplete="new-password"
                            placeholder="••••••••"
                            class="w-full h-16 bg-slate-900/40 border border-slate-800/60 rounded-xl px-6 text-base text-white placeholder-slate-600 focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500/50 transition-all hover:bg-slate-900/60 font-mono tracking-widest"
                        />
                        <div class="absolute right-6 top-1/2 -translate-y-1/2 flex items-center space-x-4">
                            <button @click="showPassword = !showPassword" type="button" class="text-slate-600 hover:text-emerald-400 transition-colors focus:outline-none p-1">
                                <Eye v-if="showPassword" class="w-5 h-5" />
                                <EyeOff v-else class="w-5 h-5" />
                            </button>
                            <div class="h-6 w-px bg-slate-800"></div>
                            <Lock class="w-5 h-5 text-slate-600 group-focus-within:text-emerald-500 transition-colors" />
                        </div>
                    </div>
                </div>

                <!-- Confirm Password -->
                <div class="space-y-4 group">
                    <label class="text-[11px] uppercase tracking-[0.3em] font-black text-slate-500 group-focus-within:text-emerald-500 transition-colors block ml-1">
                        {{ store.currentLang === 'zh' ? '确认密码' : 'Confirm Password' }}
                    </label>
                    <div class="relative">
                        <input
                            :type="showConfirmPassword ? 'text' : 'password'"
                            v-model="form.confirmPassword"
                            required
                            autocomplete="new-password"
                            placeholder="••••••••"
                            class="w-full h-16 bg-slate-900/40 border border-slate-800/60 rounded-xl px-6 text-base text-white placeholder-slate-600 focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500/50 transition-all hover:bg-slate-900/60 font-mono tracking-widest"
                        />
                        <div class="absolute right-6 top-1/2 -translate-y-1/2 flex items-center space-x-4">
                            <button @click="showConfirmPassword = !showConfirmPassword" type="button" class="text-slate-600 hover:text-emerald-400 transition-colors focus:outline-none p-1">
                                <Eye v-if="showConfirmPassword" class="w-5 h-5" />
                                <EyeOff v-else class="w-5 h-5" />
                            </button>
                            <div class="h-6 w-px bg-slate-800"></div>
                            <Lock class="w-5 h-5 text-slate-600 group-focus-within:text-emerald-500 transition-colors" />
                        </div>
                    </div>
                </div>

                <!-- Security Notice -->
                <div class="flex items-start space-x-3 p-4 rounded-2xl bg-amber-500/5 border border-amber-500/10">
                    <Shield class="w-5 h-5 text-amber-500 flex-shrink-0 mt-0.5" />
                    <p class="text-[11px] font-bold text-amber-500/80 leading-relaxed tracking-wide uppercase">
                        {{ store.currentLang === 'zh'
                            ? '这是超级管理员账户。密码不可恢复 — 请妥善保管。密钥文件 (.nanolog.key) 请务必备份。'
                            : 'This is the Super Admin account. Password cannot be recovered — store it safely. Back up your key file (.nanolog.key).'
                        }}
                    </p>
                </div>

                <!-- Submit Button -->
                <button
                    type="submit"
                    :disabled="loading"
                    class="w-full h-16 relative group overflow-hidden bg-white text-black font-black rounded-2xl shadow-[0_20px_40px_-15px_rgba(16,185,129,0.25)] transition-all transform active:scale-[0.98] disabled:opacity-50 flex items-center justify-center space-x-4"
                >
                    <div class="absolute inset-0 bg-emerald-500 transition-all duration-300 ease-out translate-y-16 group-hover:translate-y-0 opacity-10"></div>
                    <Loader2 v-if="loading" class="w-5 h-5 animate-spin" />
                    <span class="relative uppercase tracking-[0.4em] text-[13px] indent-[0.4em]">{{ loading ? (store.currentLang === 'zh' ? '初始化中...' : 'Initializing...') : t('init.submit') }}</span>
                    <ArrowRight v-if="!loading" class="w-5 h-5 group-hover:translate-x-2 transition-transform duration-500" />
                </button>
            </form>

            <footer class="mt-5 flex flex-col items-center">
                <div class="h-px w-12 bg-gradient-to-r from-transparent via-slate-800 to-transparent mb-10"></div>
                <p class="text-[10px] text-slate-700 font-bold uppercase tracking-[0.5em] mb-4">{{ t('auth.footer_note') }}</p>
                <div class="flex space-x-8">
                    <a href="https://github.com/coffersTech/nanolog" target="_blank" class="w-10 h-10 rounded-2xl bg-slate-950 border border-slate-900 flex items-center justify-center hover:border-slate-700 hover:bg-slate-900 transition-all group">
                        <svg class="w-5 h-5 text-slate-600 group-hover:text-white transition-colors" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.041-1.416-4.041-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg></a>
                </div>
            </footer>
        </div>
    </div>
  </div>
</template>

<style scoped>
* {
    transition-property: background-color, border-color, color, fill, stroke, opacity, box-shadow, transform;
    transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
    transition-duration: 150ms;
}

input:focus {
    box-shadow: 0 0 0 4px rgba(16, 185, 129, 0.1);
}

@keyframes glow-scan {
    0% { transform: translateY(-128px); opacity: 0; }
    10% { opacity: 1; }
    90% { opacity: 1; }
    100% { transform: translateY(100vh); opacity: 0; }
}

.animate-glow-scan {
    animation: glow-scan 4s cubic-bezier(0.4, 0, 0.2, 1) infinite;
}
</style>
