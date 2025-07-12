#!/bin/python3

# This script copies the given files to the specified destination directory. It will create the destination directory if it does not exist. It will also change the package import to the new location.
# example: ./copy_files.py --file ./pkg/engine/alpha_beta_engine.go ./pkg/engine/alpha_beta_engine_test.go ./pkg/engine/alpha_beta_engine_benchmark_test.go --dst-dir ./pkg/engine/new 


import os
import shutil
import argparse

def copy_files(files, dst_dir):
    """
    Copies the specified files to the destination directory.
    If the destination directory does not exist, it will be created.
    """
    print(f"Copying files to {dst_dir}...")
    for file in files:
        if not file.endswith('.go'):
            raise ValueError(f"Source file {file} is not a .go file")

    if not os.path.exists(dst_dir):
        os.makedirs(dst_dir)
    
    destination_package_name = os.path.basename(dst_dir)

    for file in files:
        src_path = file
        dst_path = os.path.join(dst_dir, os.path.basename(file))
        
        src_package_name = os.path.basename(os.path.dirname(src_path))

        # copy the file
        shutil.copy(src_path, dst_path)
        # replace the package name in the copied file
        copy_and_replace_package(dst_path, dst_path, src_package_name, destination_package_name)
        src_package_name = os.path.basename(os.path.dirname(src_path))

        # copy the file
        shutil.copy(src_path, dst_path)
        # replace the package name in the copied file
        copy_and_replace_package(dst_path, dst_path, src_package_name, destination_package_name)

def copy_and_replace_package(src_path, dst_path, old_pkg, new_pkg):
    """
    Copies a Go file and replaces the package name.
    """
    with open(src_path, 'r') as src_file:
        content = src_file.read()
    content = content.replace(f'package {old_pkg}', f'package {new_pkg}')
    with open(dst_path, 'w') as dst_file:
        dst_file.write(content)

if __name__ == "__main__":
    print("Script started", flush=True)
    print("Starting the file copy process...")
    parser = argparse.ArgumentParser(description='Copy Go files and change package imports.')
    parser.add_argument('--files', nargs='+', required=True, help='List of files to copy in the format')
    parser.add_argument('--dst-dir', required=True, help='Destination directory to copy files to')

    args = parser.parse_args()

    copy_files(args.files, args.dst_dir)
