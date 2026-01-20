# Troubleshooting

Common issues and solutions for ELMOS.

## Build Issues

### "gmake not found"
```bash
brew install make
# Use gmake instead of make
```

### Toolchain Build Fails
- Run `./build/elmos doctor` for missing deps
- Check Xcode CLT: `xcode-select --install`
- Clean and retry: `./build/elmos toolchains clean`

### Kernel Build Errors
- Regenerate config: `./build/elmos kernel config defconfig`
- Apply patches: `./build/elmos patch apply`
- Check HOSTCFLAGS: ELMOS sets automatically

## Runtime Issues

### QEMU Won't Start
- Ensure kernel and rootfs exist
- Check architecture match: `./build/elmos arch`

### Module Load Fails
- Verify toolchain compatibility
- Check kernel symbols: `modinfo module.ko`

### TUI Shows Help Text
- Rebuild ELMOS: `task build`

## Environment

### UUID Conflicts
- Patches applied? Check `./build/elmos status`

### Slow Builds
- Increase jobs: Set `JOBS=8` in config
- Use SSD for workspace

## General

### Doctor Fails
- Reinstall deps
- Update macOS to Sequoia+

### Sparse Image Issues
- Recreate: `rm build/elmos.sparseimage; ./build/elmos init`

## Getting Help

- [GitHub Issues](https://github.com/NguyenTrongPhuc552003/elmos/issues)
- Check logs in `build/`