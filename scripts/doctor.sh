#!/bin/bash
# scripts/doctor.sh
# Modular environment and dependency checker for native kernel builds + disk-based rootfs.

# ─────────────────────────────────────────────────────────────
# Helper: Colored status check
# ─────────────────────────────────────────────────────────────
check_item() {
	local desc="$1"
	shift
	if "$@"; then
		printf "  [${GREEN}OK${NC}] %s\n" "$desc"
		return 0
	else
		printf "  [${RED}FAIL${NC}] %s\n" "$desc"
		return 1
	fi
}

# ─────────────────────────────────────────────────────────────
# Helper: Taps MacOS check
# ─────────────────────────────────────────────────────────────
taps_check() {
	echo "Checking Homebrew taps:"

	local REQUIRED_TAPS=(
		messense/macos-cross-toolchains
	)

	local INSTALLED_TAPS
	INSTALLED_TAPS=$(brew tap 2>/dev/null)

	local missing=0

	for tap in "${REQUIRED_TAPS[@]}"; do
		if echo "$INSTALLED_TAPS" | grep -Fxq "$tap"; then
			check_item "Tap: $tap" true
		else
			check_item "Tap: $tap" false
			echo -e "${YELLOW}Fix:${NC} brew tap $tap"
			((issues_found++))
			missing=1
		fi
	done

	echo
	return $missing
}

# ─────────────────────────────────────────────────────────────
# 1. Check Homebrew packages
# ─────────────────────────────────────────────────────────────
check_homebrew_packages() {
	# Check required taps first
	if ! taps_check; then
		echo -e "${RED}Error:${NC} Missing required Homebrew taps. Aborting package check."
		return 1
	fi

	echo "Checking required Homebrew packages:"
	local REQUIRED=(
		llvm
		lld
		gnu-sed
		make
		libelf
		git
		qemu
		fakeroot
		e2fsprogs
		wget
		riscv64-elf-gdb
		arm-none-eabi-gdb
		aarch64-elf-gdb
	)

	# Get installed formulae (one per line)
	local INSTALLED
	INSTALLED=$(brew list --formulae 2>/dev/null | tr ' ' '\n')

	local missing=()

	for pkg in "${REQUIRED[@]}"; do
		if echo "$INSTALLED" | grep -Fxq "$pkg"; then
			check_item "Package: $pkg" true
		else
			check_item "Package: $pkg" false
			missing+=("$pkg")
			((issues_found++))
		fi
	done

	if [ ${#missing[@]} -gt 0 ]; then
		echo
		echo -e "${YELLOW}Fix:${NC} brew install ${missing[*]}"
	fi
	echo
}

# ─────────────────────────────────────────────────────────────
# 2. Check custom headers directory and shims (Dynamic Tree)
# ─────────────────────────────────────────────────────────────
check_custom_headers() {
	local headers_dir="${SCRIPT_DIR}/libraries"
	echo "Checking custom macOS kernel headers:"

	if [ ! -d "$headers_dir" ]; then
		echo -e "  [${RED}FAIL${NC}] Directory missing: $headers_dir"
		((issues_found++))
		return 1
	fi

	# 1. We run tree twice: once for visual display, once for absolute paths.
	# 2. 'sed $d' is used to delete the last line (the summary) safely on macOS.
	local tree_visual
	local tree_paths
	tree_visual=$(tree -n "$headers_dir" | sed '$d')
	tree_paths=$(tree -nf "$headers_dir" | sed '$d')

	# 3. Use a while loop to read both streams line by line
	while IFS= read -r v_line <&3 && IFS= read -r p_line <&4; do

		# Handle the root line (usually just 'libraries')
		if [[ "$v_line" == "libraries" ]]; then
			echo -e "  [${GREEN}OK${NC}] $v_line"
			continue
		fi

		# Extract the absolute path from the path line (last field)
		# Strip symlink arrows to get the actual link path for validation
		local full_path=$(echo "$p_line" | awk '{print $NF}' | sed 's|->.*$||')

		local status="[${GREEN}OK${NC}]"

		# 4. Validation Logic
		if [ -L "$full_path" ]; then
			# Check if the symlink target is actually reachable (exists)
			if [ ! -e "$full_path" ]; then
				status="[${RED}ER${NC}]"
				((issues_found++))
			fi
		elif [ ! -e "$full_path" ]; then
			status="[${RED}ER${NC}]"
			((issues_found++))
		fi

		# 5. Output with cleaned-up visual line (removing absolute path prefix)
		# Use | as delimiter in sed to handle slashes in $headers_dir
		local clean_v_line=$(echo "$v_line" | sed "s|$headers_dir/||g")
		echo -e "  $status $v_line"

	done 3<<<"$tree_visual" 4<<<"$tree_paths"
	echo
}

# ─────────────────────────────────────────────────────────────
# 3. Auto-fix missing elf.h
# ─────────────────────────────────────────────────────────────
fix_missing_elf_h() {
	local headers_dir="${MACOS_HEADERS:-$HOME/Documents/kernel-dev/linux/libraries}"

	if [ -f "$headers_dir/elf.h" ]; then
		return
	fi

	echo -e "${YELLOW}elf.h missing — download from glibc? (Y/n)${NC}"
	read -r choice
	if [[ "$choice" =~ ^[Nn]$ ]]; then
		return
	fi

	local glibc_ver="2.42"
	local url="https://raw.githubusercontent.com/bminor/glibc/glibc-${glibc_ver}/elf/elf.h"

	mkdir -p "$headers_dir"
	echo "Downloading elf.h..."
	if wget -q --show-progress "$url" -O "$headers_dir/elf.h"; then
		echo -e "  [${GREEN}FIXED${NC}] elf.h downloaded."
	else
		echo -e "  [${RED}ERROR${NC}] Download failed."
		((issues_found++))
	fi
	echo
}

# ─────────────────────────────────────────────────────────────
# 4. Check debootstrap submodule
# ─────────────────────────────────────────────────────────────
check_debootstrap() {
	echo "Checking debootstrap tool:"
	if [ -f "${DEBOOTSTRAP_PATH}/debootstrap" ]; then
		check_item "debootstrap script present" true
	else
		echo -e "  [${RED}MISSING${NC}] Not found at: ${DEBOOTSTRAP_PATH}/debootstrap"
		echo "  Fix: git submodule update --init --recursive tools/debootstrap"
		((issues_found++))
	fi
	echo
}

# ─────────────────────────────────────────────────────────────
# Main doctor function
# ─────────────────────────────────────────────────────────────
run_doctor() {
	echo -e "  [${YELLOW}INFO${NC}] Running doctor — environment check..."
	echo

	local issues_found=0

	check_homebrew_packages
	check_custom_headers
	check_debootstrap
	fix_missing_elf_h # Run after header check

	# ─────────────────────────────────────────────────────────────
	# Final verdict
	# ─────────────────────────────────────────────────────────────
	if [ "$issues_found" -eq 0 ]; then
		echo -e "  ${GREEN}All checks passed! Environment ready for build and rootfs creation.${NC}"
		return 0
	else
		echo -e "  ${RED}Found $issues_found issue(s). Please fix before proceeding.${NC}"
		return 1
	fi
}
