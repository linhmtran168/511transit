import { fromUnixTime, format } from 'date-fns'

export default function useUtils() {
  const formatUnixTime = (time: number | undefined | null, timeFormat?: string): string => {
    if (!time) {
      return '???'
    }
    if (!timeFormat) {
      timeFormat = 'yyyy/MM/dd p'
    }
    return format(fromUnixTime(time), 'yyyy/MM/dd p')
  }

  return {
    formatUnixTime
  }
}