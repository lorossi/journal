# this script automatically builds and zips the package for multiple
# architectures and operating systems.
# it is really basic, nothing too fancu.
# (almost) PEP8 compliant BTW
import os
import json
import zipfile
import subprocess


def build(binaries_folder="../binaries/", settings_file="buildsettings.json", compressed_filename="builds.zip"):
    try:
        os.remove(binaries_folder + compressed_filename)
    except FileNotFoundError:
        pass
    except Exception as e:
        print(e)

    with open(settings_file) as f:
        settings = json.load(f)

    for b in settings["builds"]:
        if not os.path.exists(binaries_folder + b["folder"]):
            os.makedirs(binaries_folder + b["folder"])
        options = f"env GOOS={b['os']} GOARCH={b['architecture']} go build " \
                  f"-o {binaries_folder}{b['folder']}{b['filename']}"

        subprocess.run(options.split(" "))

    zipf = zipfile.ZipFile(binaries_folder + compressed_filename, "w",
                           zipfile.ZIP_DEFLATED)
    for root, dirs, files in os.walk(binaries_folder):
        for file in files:
            zipf.write(os.path.join(root, file),
                       os.path.relpath(os.path.join(root, file),
                                       os.path.join(binaries_folder, '..')))
    zipf.close()


if __name__ == "__main__":
    build()
