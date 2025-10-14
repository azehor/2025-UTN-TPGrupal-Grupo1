export function isStaffAuthenticated(): boolean {
  return sessionStorage.getItem('isStaff') === '1'
}

export function logoutStaff() {
  sessionStorage.removeItem('isStaff')
}
