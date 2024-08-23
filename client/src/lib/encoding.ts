export const uuid2bytes = (uuid: string) => {
  const hex = uuid.replace(/-/g, '')

  const bytes = new Uint8Array(16)
  for (let i = 0; i < 16; i++) {
    bytes[i] = parseInt(hex.slice(i * 2, i * 2 + 2), 16)
  }

  return bytes
}

export const b64urlDecode = (input: string) => {
  return Uint8Array.from(atob(input.replace(/-/g, '+').replace(/_/g, '/')), c =>
    c.charCodeAt(0)
  )
}

export const b64urlEncode = (input: ArrayBuffer) => {
  const bytes = input instanceof Uint8Array ? input : new Uint8Array(input)

  return btoa(String.fromCharCode(...bytes))
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=/g, '')
}
