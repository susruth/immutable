package record_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/susruth/immutable/record"
)

var _ = Describe("Immutable Record", func() {
	Context("when appending elements to a record", func() {
		It("should create new versions of the record when using ints", func() {
			record := New()
			record.BatchUpdate(Snapshot{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
			})

			Expect(record.Snapshot(1)).Should(Equal(Snapshot{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
			}))
			record.Update("e", 5)
			Expect(record.Snapshot(2)).Should(Equal(Snapshot{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
				"e": 5,
			}))
		})

		It("should create new versions of the record when using strings", func() {
			record := New()
			record.BatchUpdate(Snapshot{
				"a": "hi",
				"b": "hello",
				"c": "bye",
			})

			Expect(record.Snapshot(1)).Should(Equal(Snapshot{
				"a": "hi",
				"b": "hello",
				"c": "bye",
			}))
			record.Update("d", "good")

			Expect(record.Snapshot(2)).Should(Equal(Snapshot{
				"a": "hi",
				"b": "hello",
				"c": "bye",
				"d": "good",
			}))
		})
	})

	Context("when updating elements on a record", func() {
		It("should create new versions of the record when using ints", func() {
			record := New()
			record.BatchUpdate(Snapshot{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
			})

			Expect(record.Snapshot(1)).Should(Equal(Snapshot{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
			}))

			record.Update("c", 5)
			Expect(record.Snapshot(2)).Should(Equal(Snapshot{
				"a": 1,
				"b": 2,
				"c": 5,
				"d": 4,
			}))
		})

		It("should create new versions of the record when using strings", func() {
			record := New()
			record.BatchUpdate(Snapshot{
				"a": "hi",
				"b": "hello",
				"c": "bye",
			})

			Expect(record.Snapshot(1)).Should(Equal(Snapshot{
				"a": "hi",
				"b": "hello",
				"c": "bye",
			}))
			record.Update("c", "good")

			Expect(record.Snapshot(2)).Should(Equal(Snapshot{
				"a": "hi",
				"b": "hello",
				"c": "good",
			}))
		})
	})

	Context("when removing elements from a record", func() {
		It("should create new versions of the record when using ints", func() {
			record := New()
			record.BatchUpdate(Snapshot{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
			})

			Expect(record.Snapshot(1)).Should(Equal(Snapshot{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
			}))

			record.Delete("c")
			Expect(record.Snapshot(2)).Should(Equal(Snapshot{
				"a": 1,
				"b": 2,
				"d": 4,
			}))
		})

		It("should create new versions of the record when using strings", func() {
			record := New()
			record.BatchUpdate(Snapshot{
				"a": "hi",
				"b": "hello",
				"c": "bye",
			})

			Expect(record.Snapshot(1)).Should(Equal(Snapshot{
				"a": "hi",
				"b": "hello",
				"c": "bye",
			}))
			record.Delete("c")

			Expect(record.Snapshot(2)).Should(Equal(Snapshot{
				"a": "hi",
				"b": "hello",
			}))
		})
	})

	Context("when marshalling and unmarshalling a record", func() {
		It("should create new versions of the record when using ints", func() {
			record := New()
			record.BatchUpdate(Snapshot{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
			})

			data, err := record.MarshalBinary()
			Expect(err).Should(BeNil())
			secRecord := New()
			Expect(secRecord.UnmarshalBinary(data)).Should(BeNil())

			Expect(record.Snapshot(1)).Should(Equal(Snapshot{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
			}))
		})

		It("should create new versions of the record when using strings", func() {
			record := New()
			record.BatchUpdate(Snapshot{
				"a": "hi",
				"b": "hello",
				"c": "bye",
			})

			data, err := record.MarshalBinary()
			Expect(err).Should(BeNil())
			secRecord := New()
			Expect(secRecord.UnmarshalBinary(data)).Should(BeNil())

			Expect(record.Snapshot(1)).Should(Equal(Snapshot{
				"a": "hi",
				"b": "hello",
				"c": "bye",
			}))
		})
	})
})
