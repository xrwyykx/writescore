<template>
  <form @submit.prevent="handleRegister">
    <div class="mb-4">
      <label for="register-username" class="block text-sm font-medium text-gray-700 mb-1">用户名</label>
      <input 
        id="register-username" 
        v-model="username" 
        type="text" 
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-primary focus:border-primary"
        required
      />
    </div>
    <div class="mb-4">
      <label for="register-password" class="block text-sm font-medium text-gray-700 mb-1">密码</label>
      <input 
        id="register-password" 
        v-model="password" 
        type="password" 
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-primary focus:border-primary"
        required
      />
    </div>
    <div class="mb-6">
      <label for="register-confirm-password" class="block text-sm font-medium text-gray-700 mb-1">确认密码</label>
      <input 
        id="register-confirm-password" 
        v-model="confirmPassword" 
        type="password" 
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-primary focus:border-primary"
        required
      />
    </div>
    <div class="mb-6">
      <label for="register-avatar" class="block text-sm font-medium text-gray-700 mb-1">头像</label>
      <input 
        id="register-avatar" 
        v-model="avatar" 
        type="avatar" 
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-primary focus:border-primary"
        required
      />
    </div>
    <div class="mb-6">
      <label for="register-nick_name" class="block text-sm font-medium text-gray-700 mb-1">昵称</label>
      <input 
        id="register-nick_name" 
        v-model="nickName" 
        type="nickName" 
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-primary focus:border-primary"
        
      />
    </div>
    <button 
      type="submit" 
      class="w-full bg-primary text-white py-2 px-4 rounded-md hover:bg-primary-dark focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
      :disabled="isLoading"
    >
      {{ isLoading ? '注册中...' : '注册' }}
    </button>
  </form>
</template>

<script setup>
import { ref } from 'vue';

const emit = defineEmits(['register-success']);

const username = ref('');
const password = ref('');
const confirmPassword = ref('');
const avatar = ref('');
const nickName = ref('');
const isLoading = ref(false);

const handleRegister = async () => {
  if (isLoading.value) return;
  
  // Validate passwords match
  if (password.value !== confirmPassword.value) {
    alert('两次输入的密码不一致');
    return;
  }
  
  isLoading.value = true;
  
  try {
    const response = await fetch('http://47.113.178.27:8099/api/common/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        username: username.value,
        password: password.value,
        avatar: avatar.value,
      })
    });
    
    if (!response.ok) {
      throw new Error('注册失败');
    }
    emit('register-success', username.value);
    username.value = '';
    password.value = '';
    confirmPassword.value = '';
    avatar.value = '';
    
    alert('注册成功，请登录');
  } catch (error) {
    alert('注册失败: ' + error.message);
  } finally {
    isLoading.value = false;
  }
};
</script>