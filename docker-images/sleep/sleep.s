global _start

section .data
  timespec:
    tv_sec  dq 65535
    tv_nsec dq 0

section .text
  _start:
.loop:
    mov rax, 35
    mov rdi, timespec
    xor rsi, rsi
    syscall
    jmp .loop
