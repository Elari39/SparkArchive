import { customRef, type Ref } from 'vue'

/**
 * 防抖 ref：写入后延迟 delay 毫秒才把新值同步给下游（watch / computed）。
 * 用于筛选输入框，避免逐字符触发请求。
 *
 * 注意：读取到的始终是**已提交**的值；输入框绑定的是这个 ref，
 * 因此输入响应不受影响，只是「值变化」被延迟通知。
 */
export function useDebouncedRef<T>(initial: T, delay = 280): Ref<T> {
  let value = initial
  let timer: ReturnType<typeof setTimeout> | undefined
  return customRef<T>((track, trigger) => ({
    get() {
      track()
      return value
    },
    set(next: T) {
      if (timer) clearTimeout(timer)
      timer = setTimeout(() => {
        value = next
        trigger()
      }, delay)
    },
  }))
}
