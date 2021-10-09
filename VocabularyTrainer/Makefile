buildall:
	fyne package -os linux -icon resources/icon.png
	fyne-cross darwin -icon resources/icon.png
	fyne-cross windows -icon resources/icon.png

	# turn macOS .app to dmg
	dd if=/dev/zero of=/tmp/VocabularyTrainer.dmg bs=1M count=16 status=progress 
	mkfs.hfsplus -v VocabularyTrainer /tmp/VocabularyTrainer.dmg

	sudo mkdir -pv /mnt/tmp
	sudo mount -o loop /tmp/VocabularyTrainer.dmg /mnt/tmp
	sudo cp -av fyne-cross/dist/darwin-amd64/VocabularyTrainer.app /mnt/tmp

	sudo umount /mnt/tmp

	cp /tmp/VocabularyTrainer.dmg fyne-cross/dist/darwin-amd64/VocabularyTrainer.dmg

clean:
	if [ -d "fyne-cross" ]; then rm -r fyne-cross; fi
	if [ -f "VocabularyTrainer" ]; then rm -r VocabularyTrainer; fi
