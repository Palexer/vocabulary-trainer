buildall:
	@echo "Creating dist/ directory"
	@mkdir -p dist/win
	@mkdir -p dist/mac
	@mkdir -p dist/linux
	
	@echo "Packaging for Linux"
	@fyne package -os linux -icon resources/icon.png -appID io.github.palexer.vocabulary-trainer -release -appVersion "1.3.1"
	@echo "Packaging for Darwin"
	@fyne package -os darwin -icon resources/icon.png -appID io.github.palexer.vocabulary-trainer -release -appVersion "1.3.1"
	@echo "Packaging for Windows"
	@#you have to install the x86_64-w64-mingw64-gcc compiler to cross compile for Windows, see https://developer.fyne.io/started/cross-compiling for more
	@CC=x86_64-w64-mingw32-gcc fyne package -os windows -icon resources/icon.png -appID io.github.palexer.vocabulary-trainer -release -appVersion "1.3.1"

	# turn macOS .app to dmg
	@dd if=/dev/zero of=/tmp/VocabularyTrainer.dmg bs=1M count=16 status=progress 
	@mkfs.hfsplus -v VocabularyTrainer /tmp/VocabularyTrainer.dmg

	@sudo mkdir -pv /mnt/tmp
	@sudo mount -o loop /tmp/VocabularyTrainer.dmg /mnt/tmp
	@sudo cp -av ./VocabularyTrainer.app /mnt/tmp

	@sudo umount /mnt/tmp

	@cp /tmp/VocabularyTrainer.dmg dist/mac/VocabularyTrainer.dmg
	@mv ./VocabularyTrainer.tar.gz dist/linux/VocabularyTrainer.tar.gz
	@mv ./VocabularyTrainer.exe dist/win/VocabularyTrainer.exe
	@rm -r VocabularyTrainer.app
	@rm VocabularyTrainer
	@echo "Finished."

clean:
	@if [ -d "fyne-cross" ]; then rm -r fyne-cross; fi
	@if [ -d "dist" ]; then rm -r dist; fi
	@if [ -d "VocabularyTrainer.app" ]; then rm -r VocabularyTrainer.app; fi
	@if [ -f "VocabularyTrainer" ]; then rm -r VocabularyTrainer; fi
	@if [ -f "VocabularyTrainer.exe" ]; then rm -r VocabularyTrainer.exe; fi
	@if [ -f "VocabularyTrainer.tar.gz" ]; then rm -r VocabularyTrainer.tar.gz; fi
