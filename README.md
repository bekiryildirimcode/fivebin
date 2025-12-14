# Vocabulary Learner

Modern ve hafıza verimli bir kelime öğrenme uygulaması. Go ile geliştirilmiş, macOS'ta çalışır.

> 💡 **Not**: Bu uygulama Cursor AI asistanı kullanılarak kod yazmadan, sadece doğal dil komutlarıyla geliştirilmiştir.

## Özellikler

- 🎯 **Hafıza Verimli**: Çok büyük sözlük dosyalarını RAM'e yüklemeden işler
- 🔍 **Kelime Arama**: Hızlı arama özelliği
- 🎲 **Rastgele Kelime**: Tekrarsız rastgele kelime öğrenme
- ✏️ **Kişisel Notlar**: Her kelime için kendi notlarınızı kaydedin
- 💾 **Otomatik Kayıt**: Notlarınız otomatik olarak kaydedilir
- 🎨 **Modern Arayüz**: Temiz ve kullanıcı dostu tasarım

## Kullanım

### macOS'ta Çalıştırma

1. `Vocabulary Learner.app` dosyasını çift tıklayın
2. Uygulama otomatik olarak açılır (sözlük dosyası içine gömülü)

### Temel İşlemler

- **Rastgele Kelime**: "🎲 New Word" butonuna tıklayın
- **Kelime Ara**: Üstteki arama kutusuna kelime yazın ve Enter'a basın
- **Not Kaydet**: Kelime için notlarınızı yazın ve "💾 Save Meaning" butonuna tıklayın

### Görüntülenen Bilgiler

- Kelime (büyük, ortalanmış)
- Kelime türü (isim, fiil, vb.)
- Seviye (A1, B2, vb.)
- Telaffuz (US / UK)
- Örnek cümleler

## Teknik Bilgi

- **Dil**: Go
- **GUI**: Fyne
- **Sözlük**: Uygulamaya gömülü (resources/data.json)
- **Notlar**: user_meanings.json dosyasında saklanır

## Geliştirme

Bu uygulama **Cursor AI** ile geliştirilmiştir:
- Kod yazılmadan, sadece doğal dil komutlarıyla
- Streaming JSON parsing ile hafıza verimli çalışır
- Tüm özellikler AI asistanı tarafından oluşturuldu

## Sorun Giderme

- **Uygulama açılmıyor**: macOS 11.0 veya üzeri gerekli
- **Arama çalışmıyor**: Kelime tam eşleşme gerektirir (büyük/küçük harf duyarsız)

---

**İyi öğrenmeler! 📚✨**
