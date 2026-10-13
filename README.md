# moa

## Opcodes
| Opcode | Byte | Params | Description |
|--------|------|--------|-------------|
| NOP    | 0x00 | | No op |
| HLT    | 0x01 | | Halts execution |
| LDI    | 0x10 | register, value | Loads instantaneous value into register |
| STI    | 0x11 | location, value | Stores instantaneous value in memory |
| LD     | 0x12 | register, location | Loads value from memory into register |
| ST     | 0x13 | location, register | Stores value from register in memory |
| MOV    | 0x14 | destination, source | Copies value from source register to destination register |

## Todo
- [ ] Proper error handling

## File formats
- `.moa` - moa assembly source (not implemented)
- `.mxe` - compiled moa bytecode (not implemented)
- `.mhrbc` - moa human readable bytecode
