import { ref } from 'vue';

export default function useLoading(initValue = false) {
  const loading = ref(initValue);

  function startLoading () {
    loading.value = true;
  }

  function endLoading () {
    loading.value = false;
  }


  return {
    loading,
    startLoading,
    endLoading
  };
}
